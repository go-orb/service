// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/go-orb/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package service

import (
	"context"
	"fmt"
	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/server"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/service/httpgateway/pkg/httpgateway_handler"
	"github.com/go-orb/service/httpgateway/pkg/httpgateway_server"
)

// Injectors from wire.go:

// ProvideRunner provides a runner for the service.
func ProvideRunner(serviceContext *cli.ServiceContext, appconfigData cli.AppConfigData, flags []*cli.Flag) (Runner, error) {
	v, err := types.ProvideComponents()
	if err != nil {
		return nil, err
	}
	serviceContextHasConfigData, err := cli.ProvideServiceConfigData(serviceContext, appconfigData, flags)
	if err != nil {
		return nil, err
	}
	logger, err := log.ProvideWithServiceNameField(serviceContextHasConfigData, serviceContext, v)
	if err != nil {
		return nil, err
	}
	registryType, err := registry.ProvideNoOpts(serviceContext, v, logger)
	if err != nil {
		return nil, err
	}
	clientType, err := client.ProvideNoOpts(serviceContext, v, logger, registryType)
	if err != nil {
		return nil, err
	}
	httpgateway_serverServer, err := httpgateway_server.Provide(serviceContext, v, logger, clientType)
	if err != nil {
		return nil, err
	}
	serverServer, err := server.ProvideNoOpts(serviceContext, v, logger, registryType)
	if err != nil {
		return nil, err
	}
	handler, err := httpgateway_handler.Provide(logger, serverServer, v, httpgateway_serverServer)
	if err != nil {
		return nil, err
	}
	actionServer, err := provideActionServer(serviceContext, v, logger, httpgateway_serverServer, handler)
	if err != nil {
		return nil, err
	}
	actionHealth, err := provideActionHealth()
	if err != nil {
		return nil, err
	}
	runner, err := provideRunner(actionServer, actionHealth)
	if err != nil {
		return nil, err
	}
	return runner, nil
}

// wire.go:

type Runner func(command []string) error

type ActionServer func() error

type ActionHealth func() error

func provideActionServer(
	serviceContext *cli.ServiceContext,
	components *types.Components,
	logger log.Logger, server2 *httpgateway_server.Server,
	handler *httpgateway_handler.Handler,
) (ActionServer, error) {
	return func() error {
		if err := components.Add(server2, types.PriorityHandler-10); err != nil {
			return err
		}
		if err := components.Add(handler, types.PriorityHandler); err != nil {
			return err
		}

		for _, c := range components.Iterate(false) {
			logger.Debug("Starting", "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))

			err := c.Start(serviceContext.Context())
			if err != nil {
				logger.Error("Failed to start", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
				return fmt.Errorf("failed to start component %s/%s: %w", c.Type(), c.String(), err)
			}
		}

		<-serviceContext.Context().Done()

		ctx := context.Background()

		for _, c := range components.Iterate(true) {
			logger.Debug("Stopping", "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))

			err := c.Stop(ctx)
			if err != nil {
				logger.Error("Failed to stop", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
			}
		}

		serviceContext.StopWaitGroup().Done()

		return nil
	}, nil
}

func provideActionHealth() (ActionHealth, error) {
	return func() error {
		return nil
	}, nil
}

func provideRunner(actionServer ActionServer, actionHealth ActionHealth) (Runner, error) {
	return func(command []string) error {
		switch command[0] {
		case "server":
			return actionServer()
		case "health":
			return actionHealth()
		default:
			return fmt.Errorf("unknown action: %s", command[0])
		}
	}, nil
}
