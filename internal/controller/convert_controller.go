package controller

import (
	"go-img-processing/bootstrap"
	"go-img-processing/internal/service"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gocv.io/x/gocv"
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

	if !strings.HasSuffix(strings.ToLower(file.Filename), string(gocv.PNGFileExt)) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "only PNG files are accepted"})
		return
	}

	src, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	outputFile, res, err := c.service.Convert(ctx, src)
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
