// Package httpgateway_server provides the http gateway.
package httpgateway_server

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"sync"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/codecs"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/orberrors"
	"github.com/go-orb/service/httpgateway/proto/httpgateway_v1"
)

var _ types.Component = (*Server)(nil)

// Server is a component implementing the http gateway.
type Server struct {
	config Config

	logger log.Logger
	client client.Type

	server *http.Server

	mu     sync.Mutex
	routes []*httpgateway_v1.Route
}

// Type returns the type of the component.
func (s *Server) Type() string {
	return "gateway"
}

// String returns the name of the component.
func (s *Server) String() string {
	return "http"
}

// Start starts the httpgateway.
func (s *Server) Start(_ context.Context) error {
	if !s.config.Enabled {
		s.logger.Warn("Skipping start because component is disabled")
		return nil
	}

	s.server = &http.Server{
		Addr: s.config.Address,
	}

	go func() {
		s.logger.Info("Starting", "address", s.config.Address)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Failure while running", "error", err, "address", s.config.Address)
		}
	}()

	return nil
}

// Stop stops the httpgateway.
func (s *Server) Stop(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}

func (s *Server) proxyFor(route *httpgateway_v1.Route) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Map query/path params
		params := make(map[string]string)
		for _, p := range route.GetParams() {
			if len(c.Query(p)) > 0 {
				params[p] = c.Query(p)
			}
		}
		for _, p := range route.GetParams() {
			if len(c.Param(p)) > 0 {
				params[p] = c.Param(p)
			}
		}

		// Bind the request if POST/PATCH/PUT
		request := gin.H{}
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPatch || c.Request.Method == http.MethodPut {
			// Multipart form handler.
			// TODO(jochumdev): replace this with streaming.
			mf, err := c.MultipartForm()
			if err == nil {
				for k, files := range mf.File {
					for _, file := range files {
						fp, err := file.Open()
						if err != nil {
							continue
						}
						data, err := io.ReadAll(fp)
						if err != nil {
							continue
						}

						if len(files) > 1 {
							if _, ok := request[k]; !ok {
								request[k] = []string{base64.StdEncoding.EncodeToString(data)}
							} else {
								request[k] = append(request[k].([]string), base64.StdEncoding.EncodeToString(data))
							}
						} else {
							request[k] = base64.StdEncoding.EncodeToString(data)
						}
					}
				}

				for k, v := range mf.Value {
					if len(v) > 1 {
						request[k] = v
					} else {
						request[k] = v[0]
					}

				}
			} else {
				if c.ContentType() == "" {
					c.JSON(http.StatusUnsupportedMediaType, gin.H{
						"errors": []gin.H{
							{
								"code":    500,
								"message": "unsupported media-type, provide a content-type header",
							},
						},
					})
					c.Abort()
					return
				}
				if err := c.ShouldBind(&request); err != nil {
					c.Abort()
					return
				}
			}
		}

		// Set query/route params to the request
		for pn, p := range params {
			request[pn] = p
		}

		// remote request
		response, err := client.Request[gin.H](c.Request.Context(), s.client, route.GetService(), route.GetMethod(), request, client.WithContentType(codecs.MimeJSON))
		if err != nil {
			orbErr := orberrors.From(err)
			c.JSON(orbErr.Code, gin.H{
				"errors": []gin.H{
					{
						"code":    orbErr.Code,
						"message": orbErr.Error(),
					},
				},
			})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) buildGin() error {
	if s.server == nil {
		return nil
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(sloggin.New(s.logger.Logger))
	router.Use(gin.Recovery())

	for _, route := range s.routes {
		router.Handle(route.GetHttpMethod(), route.GetPath(), s.proxyFor(route))
	}

	s.server.Handler = router

	return nil
}

// Add adds the given routes to the gateway if there is no such route.
func (s *Server) Add(routes *httpgateway_v1.Routes) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, route := range routes.Routes {
		if !slices.ContainsFunc(s.routes, func(r *httpgateway_v1.Route) bool {
			return r.GetPath() == route.GetPath() && r.GetMethod() == route.GetMethod()
		}) {
			s.logger.Trace("adding route", "method", route.GetHttpMethod(), "path", route.GetPath())
			s.routes = append(s.routes, route)
		}
	}

	return s.buildGin()
}

// Set sets the given routes by path on the gateway.
func (s *Server) Set(routes *httpgateway_v1.Routes) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, route := range routes.Routes {
		s.logger.Trace("adding route", "method", route.GetHttpMethod(), "path", route.GetPath())

		found := false
		for idx, r := range s.routes {
			if r.GetPath() == route.GetPath() && r.GetMethod() == route.GetMethod() {
				found = true
				s.routes[idx] = route
			}
		}

		if !found {
			s.routes = append(s.routes, route)
		}
	}

	return s.buildGin()
}

// Remove removes the given routes from the gateway.
func (s *Server) Remove(paths *httpgateway_v1.Paths) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.routes = slices.DeleteFunc(s.routes, func(route *httpgateway_v1.Route) bool {
		return slices.ContainsFunc(paths.Paths, func(r string) bool {
			s.logger.Trace("deleting route", "method", route.GetMethod(), "path", route.GetPath())
			return route.GetPath() == r
		})
	})

	return s.buildGin()
}

// New creates a new httpgateway component.
func New(cfg Config, logger log.Logger, client client.Type) *Server {

	return &Server{
		config: cfg,
		logger: logger,
		client: client,

		routes: make([]*httpgateway_v1.Route, 0),
	}
}

// Provide provides the httpgateway component.
func Provide(
	svcCtx *cli.ServiceContext,
	components *types.Components,
	logger log.Logger,
	client client.Type,
) (*Server, error) {
	cfg := NewConfig()

	if err := config.Parse(nil, DefaultConfigSection, svcCtx.Config, &cfg); err != nil && !errors.Is(err, config.ErrNoSuchKey) {
		return nil, err
	}

	// Configure the logger.
	cLogger, err := logger.WithConfig([]string{DefaultConfigSection}, svcCtx.Config)
	if err != nil {
		return nil, err
	}

	cLogger = cLogger.With(slog.String("component", "httpgateway"))

	gateway := New(cfg, cLogger, client)

	return gateway, nil
}
