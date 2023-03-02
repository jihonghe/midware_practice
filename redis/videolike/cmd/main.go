package main

import (
	rediscli "pkg/redis"
	"videolike/internal/config"
	"videolike/internal/server"
)

func main() {
	config.Init()
	rediscli.Init(config.Config)
	server.Run()
}
