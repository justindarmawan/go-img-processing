package controller

import (
	"go-img-processing/bootstrap"
	"go-img-processing/internal/service"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gocv.io/x/gocv"
)

type CompressController struct {
	service service.CompressService
	config  *bootstrap.Container
}

func NewCompressController(server *gin.Engine, config *bootstrap.Container, service service.CompressService) {
	controller := &CompressController{
		service: service,
		config:  config,
	}

	routes := server.Group("/compress")
	{
		routes.POST("/", controller.Compress)
	}
}

func (c *CompressController) Compress(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to get image from form-data"})
		return
	}

	quality, err := strconv.Atoi(ctx.PostForm("quality"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "quality value must be number"})
		return
	}

	if quality < 0 || quality > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "quality value must between 0 and 100"})
		return
	}

	extArray := []string{string(gocv.PNGFileExt), string(gocv.JPEGFileExt), ".jpeg"}
	extension := filepath.Ext(file.Filename)
	extList := strings.Join(extArray, " ")
	if found := strings.Contains(extList, extension); !found {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "only PNG, JPG, JPEG files are accepted"})
		return
	}

	src, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	outputFile, res, err := c.service.Compress(ctx, src, quality)
	if !res {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		outputFile.Close()
		os.Remove(outputFile.Name())
	}()

	ctx.Header("Content-Disposition", "attachment; filename=compress"+ctx.PostForm("quality")+"-"+strings.Split(file.Filename, ".")[0]+".jpg")
	ctx.File(outputFile.Name())
}
