package main

import (
	"github.com/arifai/zenith/pkg/logger"
	"github.com/arifai/zenith/pkg/server"
)

func main() {
	logger.InitLogger()
	server.Run()
}
