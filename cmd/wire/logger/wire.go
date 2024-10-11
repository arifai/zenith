//go:build wireinject

package logger

import (
	"github.com/arifai/zenith/pkg/logger"
	"github.com/google/wire"
)

func ProvideLogger() logger.Logger {
	wire.Build(logger.New)
	return logger.Logger{}
}
