package controller

import (
	"go-img-processing/bootstrap"
	"os"

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
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Failed to parse form-data"})
		return
	}

	file := form.File["image"]
	if len(file) != 1 {
		ctx.JSON(400, gin.H{"error": "Invalid file upload"})
		return
	}

	src, err := file[0].Open()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	// Read the uploaded image using GoCV
	inputImage := gocv.IMRead(file[0].Filename, gocv.IMReadColor)
	if inputImage.Empty() {
		ctx.JSON(500, gin.H{"error": "Failed to read uploaded image"})
		return
	}
	defer inputImage.Close()

	// Convert PNG image to JPEG image
	outputImage := gocv.NewMat()
	defer outputImage.Close()

	// Write the converted image to a temporary file
	tempFile, err := os.CreateTemp("", "output-*.jpg")
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create temporary file"})
		return
	}
	defer tempFile.Close()

	// Write the output image to the temporary file
	if ok := gocv.IMWrite(tempFile.Name(), outputImage); !ok {
		ctx.JSON(500, gin.H{"error": "Failed to write output image"})
		return
	}

	// Return the path to the converted image
	ctx.JSON(200, gin.H{"message": "Image converted successfully", "path": tempFile.Name()})
}
