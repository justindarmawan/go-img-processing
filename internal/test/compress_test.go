package test

import (
	"context"
	"go-img-processing/bootstrap"
	"go-img-processing/internal/service"
	"io"
	"os"
	"testing"
)

func TestCompress(t *testing.T) {
	config := &bootstrap.Container{}
	compressService := service.NewCompressService(config)

	ctx := context.Background()
	imageFile, err := os.Open("../../static/img/balloon.png")
	if err != nil {
		t.Fatalf("failed to open image file: %v", err)
	}
	defer imageFile.Close()

	var reader io.Reader = imageFile
	quality := 80

	outputFile, res, err := compressService.Compress(ctx, reader, quality)

	if err != nil {
		t.Errorf("compress failed: %v", err)
	}

	if !res {
		t.Errorf("compression was not successful")
	}

	if outputFile == nil {
		t.Errorf("temporary file not created")
	} else {
		defer os.Remove(outputFile.Name())
	}
}
