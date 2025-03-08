// Package main implements the httpgateway example app
package main

import (
	"fmt"
	"os"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/service/httpgateway/pkg/service"
	"github.com/go-orb/service/httpgateway/pkg/version"

	_ "github.com/go-orb/plugins-experimental/registry/mdns"
	_ "github.com/go-orb/plugins/client/middleware/log"
	_ "github.com/go-orb/plugins/client/orb"
	_ "github.com/go-orb/plugins/client/orb_transport/grpc"
	_ "github.com/go-orb/plugins/codecs/goccyjson"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/codecs/yaml"
	_ "github.com/go-orb/plugins/config/source/file"
	_ "github.com/go-orb/plugins/log/slog"
	_ "github.com/go-orb/plugins/registry/memory"
	_ "github.com/go-orb/plugins/server/grpc"
)

func main() {
	app := cli.App{
		Name:     version.Name,
		Version:  version.Version,
		Usage:    "A foobar example app",
		NoAction: true,
		Flags: []*cli.Flag{
			{
				Name:        "log_level",
				Default:     "INFO",
				EnvVars:     []string{"LOG_LEVEL"},
				ConfigPaths: []cli.FlagConfigPath{{Path: []string{"logger", "level"}}},
				Usage:       "Set the log level, one of TRACE, DEBUG, INFO, WARN, ERROR",
			},
			{
				Name:        "address",
				Default:     ":8080",
				EnvVars:     []string{"HTTPGATEWAY_ADDRESS"},
				ConfigPaths: []cli.FlagConfigPath{{Path: []string{"httpgateway", "address"}}},
				Usage:       "Set the address to listen on",
			},
		},
		Commands: []*cli.Command{},
	}

	app.Commands = append(app.Commands, service.Commands()...)

	appContext := cli.NewAppContext(&app)

	_, err := run(appContext, os.Args)
	if err != nil {
		fmt.Printf("run error: %s\n", err)
		os.Exit(1)
	}
}
