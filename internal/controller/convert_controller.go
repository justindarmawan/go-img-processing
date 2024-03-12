package controller

import (
	"go-img-processing/bootstrap"
	"go-img-processing/internal/service"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type ConvertController struct {
	service service.ConvertService
	config  *bootstrap.Container
}

func NewConvertController(server *gin.Engine, config *bootstrap.Container, service service.ConvertService) {
	controller := &ConvertController{
		service: service,
		config:  config,
	}

	routes := server.Group("/convert")
	{
		routes.POST("/", controller.Convert)
	}
}

func (c *ConvertController) Convert(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to get image from form-data"})
		return
	}

	outputFile, res, err := c.service.Convert(ctx, file)
	if !res {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		outputFile.Close()
		os.Remove(outputFile.Name())
	}()

	ctx.Header("Content-Disposition", "attachment; filename="+strings.Split(file.Filename, ".")[0]+".jpg")
	ctx.File(outputFile.Name())
}
