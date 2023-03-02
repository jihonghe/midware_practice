package main

import (
	"videolike/internal/config"
	rediscli "videolike/internal/pkg/redis"
	"videolike/internal/server"
)

func main() {
	config.Init()
	rediscli.Init(config.Config)
	server.Run()
}
