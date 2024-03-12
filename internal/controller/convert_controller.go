package controller

import (
	"go-img-processing/bootstrap"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gocv.io/x/gocv"
)

func NewConvertController(server *gin.Engine, config *bootstrap.Container) {
	routes := server.Group("/convert")
	{
		routes.POST("/", Convert)
	}
}

func Convert(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get image from form-data"})
		return
	}

	if !strings.HasSuffix(strings.ToLower(file.Filename), ".png") {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only PNG files are accepted"})
		return
	}

	src, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	dst, err := os.CreateTemp("", "uploaded-*.png")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary file"})
		return
	}
	defer func() {
		dst.Close()
		os.Remove(dst.Name())
	}()

	_, err = io.Copy(dst, src)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy file contents"})
		return
	}

	inputImage := gocv.IMRead(dst.Name(), gocv.IMReadColor)
	if inputImage.Empty() {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read uploaded image"})
		return
	}
	defer inputImage.Close()

	tempFile, err := os.CreateTemp("", "output-*.jpg")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary file"})
		return
	}
	defer func() {
		dst.Close()
		os.Remove(tempFile.Name())
	}()

	if ok := gocv.IMWrite(tempFile.Name(), inputImage); !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write output image"})
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename="+strings.Split(file.Filename, ".")[0]+".jpg")
	ctx.File(tempFile.Name())
}
