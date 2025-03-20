//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/plugins/cli/urfave"
	"github.com/go-orb/service/httpgateway/pkg/service"
	"github.com/go-orb/service/httpgateway/pkg/version"
	"github.com/go-orb/wire"
)

func provideServiceContext(appContext *cli.AppContext) (*cli.ServiceContext, error) {
	return cli.NewServiceContext(appContext, version.Name, version.Version), nil
}

type wireRunResult struct{}

func wireRun(
	appContext *cli.AppContext,
	serviceRunner service.Runner,
) (wireRunResult, error) {
	appContext.StopWaitGroup.Add(1)
	return wireRunResult{}, serviceRunner(appContext.SelectedCommand)
}

func run(
	appContext *cli.AppContext,
	args []string,
) (wireRunResult, error) {
	panic(wire.Build(
		urfave.ProvideParser,
		cli.ProvideParsedFlagsFromArgs,
		cli.ProvideAppConfigData,

		provideServiceContext,
		service.ProvideRunner,

		wireRun,
	))
}
