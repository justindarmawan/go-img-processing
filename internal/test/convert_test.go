package test

import (
	"context"
	"go-img-processing/bootstrap"
	"go-img-processing/internal/service"
	"io"
	"os"
	"testing"
)

func TestConvert(t *testing.T) {
	config := &bootstrap.Container{}
	convertService := service.NewConvertService(config)

	ctx := context.Background()
	imageFile, err := os.Open("../../static/img/balloon.png")
	if err != nil {
		t.Fatalf("failed to open image file: %v", err)
	}
	defer imageFile.Close()

	var reader io.Reader = imageFile

	outputFile, res, err := convertService.Convert(ctx, reader)

	if err != nil {
		t.Errorf("convert failed: %v", err)
	}

	if !res {
		t.Errorf("conversion was not successful")
	}

	if outputFile == nil {
		t.Errorf("output file not created")
	} else {
		defer os.Remove(outputFile.Name())
	}
}
