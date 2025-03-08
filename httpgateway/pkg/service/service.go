// Package service implements the httpgateway service.
package service

import (
	"fmt"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/service/httpgateway/pkg/version"
)

// Name is the service name.
//
//nolint:gochecknoglobals
var Name = version.Name

// MainCommands returns the commands which get appended to the "main/monolith" App.
func MainCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "httpgateway",
			Service:     Name,
			Category:    "service",
			Subcommands: Commands(),
			NoAction:    true,
		},
	}
}

// Commands returns commands specific to the service.
func Commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "server",
			Service: Name,
			Usage:   fmt.Sprintf("Start the %s server", Name),
		},
		{
			Name:    "health",
			Service: Name,
			Usage:   fmt.Sprintf("Check the health of the %s service", Name),
		},
	}
}

// ProvideClientOpts provides options for the go-orb client.
func ProvideClientOpts() ([]client.Option, error) {
	return []client.Option{client.WithClientMiddleware(client.MiddlewareConfig{Name: "log"})}, nil
}
