package service

import (
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/pkg/logger"
)

// Service encapsulates common dependencies such as configuration and logging for use in other services.
type Service struct {
	config *config.Config
	log    logger.Logger
}

// New initializes a new Service instance with the provided configuration and logger.
func New(config *config.Config, log logger.Logger) *Service {
	return &Service{config: config, log: log}
}
