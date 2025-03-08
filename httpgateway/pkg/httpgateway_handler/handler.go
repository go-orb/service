// Package httpgateway_handler provides the http gateway handler.
package httpgateway_handler

import (
	"context"

	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/server"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/service/httpgateway/pkg/httpgateway_server"
	"github.com/go-orb/service/httpgateway/proto/httpgateway_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ types.Component = (*Handler)(nil)
var _ httpgateway_v1.HttpGatewayHandler = (*Handler)(nil)

type Handler struct {
	logger  log.Logger
	server  server.Server
	gateway *httpgateway_server.Server
}

func (h *Handler) Start(_ context.Context) error {
	hRegister := httpgateway_v1.RegisterHttpGatewayHandler(h)

	// Add our server handler to all entrypoints.
	h.server.GetEntrypoints().Range(func(_ string, entrypoint server.Entrypoint) bool {
		entrypoint.AddHandler(hRegister)

		return true
	})

	return nil
}

func (h *Handler) Stop(ctx context.Context) error {
	return nil
}

func (h *Handler) Type() string {
	return "handler"
}

func (h *Handler) String() string {
	return httpgateway_v1.HandlerHttpGateway
}

func (h *Handler) AddRoutes(ctx context.Context, routes *httpgateway_v1.Routes) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, h.gateway.Add(routes)
}

func (h *Handler) SetRoutes(ctx context.Context, routes *httpgateway_v1.Routes) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, h.gateway.Set(routes)
}

func (h *Handler) RemoveRoutes(ctx context.Context, paths *httpgateway_v1.Paths) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, h.gateway.Remove(paths)
}

func New(logger log.Logger, server server.Server, gateway *httpgateway_server.Server) *Handler {
	return &Handler{
		logger:  logger,
		server:  server,
		gateway: gateway,
	}
}

func Provide(
	logger log.Logger,
	server server.Server,
	components *types.Components,
	gateway *httpgateway_server.Server,
) (*Handler, error) {
	handler := New(logger, server, gateway)

	return handler, nil
}
