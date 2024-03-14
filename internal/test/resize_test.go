package test

import (
	"context"
	"go-img-processing/bootstrap"
	"go-img-processing/internal/service"
	"io"
	"os"
	"testing"

	"gocv.io/x/gocv"
)

func TestResize(t *testing.T) {
	config := &bootstrap.Container{}
	resizeService := service.NewResizeService(config)

	ctx := context.Background()
	imageFile, err := os.Open("../../static/img/manpencil.jpg")
	if err != nil {
		t.Fatalf("failed to open image file: %v", err)
	}
	defer imageFile.Close()

	var reader io.Reader = imageFile
	width := 300
	height := 300
	ar := true
	scalar := gocv.NewScalar(255, 0, 0, 255)
	extension := ".jpg"

	outputFile, res, err := resizeService.Resize(ctx, reader, width, height, ar, &scalar, extension)

	if err != nil {
		t.Errorf("resize failed: %v", err)
	}

	if !res {
		t.Errorf("resizing was not successful")
	}

	if outputFile == nil {
		t.Errorf("output file not created")
	} else {
		defer os.Remove(outputFile.Name())
	}
}
