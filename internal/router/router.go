// internal/router/router.go
package router

import (
	"myapi/internal/handler"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	r.GET("/health", handler.Health)

	v1 := r.Group("/v1")
	{
		v1.GET("/users/:id", handler.GetUser)
		v1.POST("/users", handler.CreateUser)
	}

	return r
}
