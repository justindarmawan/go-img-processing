package controller

import (
	"go-img-processing/bootstrap"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHealthController(server *gin.Engine, config *bootstrap.Container) {
	routes := server.Group("/health")
	{
		routes.POST("/", Health)
		routes.HEAD("/", Health)
		routes.GET("/", Health)
	}
}

func Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ".")
}
