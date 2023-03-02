package video

import (
	"github.com/gin-gonic/gin"
)

func Run(r gin.IRouter) {
	registerRouters(r)
}

func registerRouters(r gin.IRouter) {
	var v1 = r.Group("/video/v1")
	{
		v1.POST("/like", like)
		v1.POST("/unlike", unlike)
	}
}
