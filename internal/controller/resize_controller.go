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

type ResizeController struct {
	service service.ResizeService
	config  *bootstrap.Container
}

func NewResizeController(server *gin.Engine, config *bootstrap.Container, service service.ResizeService) {
	controller := &ResizeController{
		service: service,
		config:  config,
	}

	routes := server.Group("/resize")
	{
		routes.POST("/", controller.Resize)
	}
}

func (c *ResizeController) Resize(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to get image from form-data"})
		return
	}

	width, err := strconv.Atoi(ctx.PostForm("width"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "width value must be number"})
		return
	}

	if width <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "width value must be greater than 0"})
		return
	}

	height, err := strconv.Atoi(ctx.PostForm("height"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "height value must be number"})
		return
	}

	if width*height > 7680*4320 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "pixel more than 8K (7680*4320)"})
		return
	}

	if height <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "height value must be greater than 0"})
		return
	}

	ar := false
	if ctx.PostForm("lockAspectRatio") != "" {
		ar, err = strconv.ParseBool(ctx.PostForm("lockAspectRatio"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "lockAspectRatio value must be true or false"})
			return
		}
	}

	var scalar *gocv.Scalar
	bgcolor := ctx.PostForm("bgColorRGB")
	if bgcolor != "" {
		bgred := strings.Split(bgcolor, ",")[0]
		red, err := strconv.Atoi(bgred)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "red value must be number"})
			return
		}

		bggreen := strings.Split(bgcolor, ",")[1]
		green, err := strconv.Atoi(bggreen)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "green value must be number"})
			return
		}

		bgblue := strings.Split(bgcolor, ",")[2]
		blue, err := strconv.Atoi(bgblue)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "blue value must be number"})
			return
		}
		col := gocv.NewScalar(float64(blue), float64(green), float64(red), 255)
		scalar = &col
	} else {
		scalar = nil
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

	outputFile, res, err := c.service.Resize(ctx, src, width, height, ar, scalar, extension)
	if !res {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		outputFile.Close()
		os.Remove(outputFile.Name())
	}()

	ctx.Header("Content-Disposition", "attachment; filename=resize-"+file.Filename)
	ctx.File(outputFile.Name())
}
