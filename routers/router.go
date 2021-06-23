package routers

import (
	"github.com/chenobdo/go-gin-example/pkg/setting"
	"github.com/gin-gonic/gin"
)

func InitRoute() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping success",
		})
	})

	return r
}