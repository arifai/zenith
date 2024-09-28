//go:build wireinject

package config

import (
	"github.com/arifai/zenith/config"
	"github.com/google/wire"
)

func ProvideConfig(filenames ...string) *config.Config {
	wire.Build(config.NewConfig)
	return nil
}
