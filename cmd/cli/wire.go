//go:build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/tnb-labs/panel/internal/app"
	"github.com/tnb-labs/panel/internal/apps"
	"github.com/tnb-labs/panel/internal/bootstrap"
	"github.com/tnb-labs/panel/internal/data"
	"github.com/tnb-labs/panel/internal/route"
	"github.com/tnb-labs/panel/internal/service"
)

// initCli init command line.
func initCli() (*app.Cli, error) {
	panic(wire.Build(bootstrap.ProviderSet, route.ProviderSet, service.ProviderSet, data.ProviderSet, apps.ProviderSet, app.NewCli))
}
