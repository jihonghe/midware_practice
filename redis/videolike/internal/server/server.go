package server

import (
	"github.com/gin-gonic/gin"

	"videolike/internal/apps/video"
	"videolike/internal/config"
)

func Run() {
	var engin = gin.Default()
	video.Run(engin.Group("/"))

	err := engin.Run(config.Config.GetString("http.addr"))
	if err != nil {
		panic(err)
	}
}
