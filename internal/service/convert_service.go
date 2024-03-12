package service

import (
	"context"
	"errors"
	"go-img-processing/bootstrap"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"gocv.io/x/gocv"
)

type ConvertService interface {
	Convert(ctx context.Context, src *multipart.FileHeader) (*os.File, bool, error)
}

type ConvertServiceImpl struct {
	cfg *bootstrap.Container
}

func NewConvertService(cfg *bootstrap.Container) ConvertService {
	return &ConvertServiceImpl{cfg: cfg}
}

func (s *ConvertServiceImpl) Convert(ctx context.Context, file *multipart.FileHeader) (*os.File, bool, error) {
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".png") {
		return nil, false, errors.New("only PNG files are accepted")
	}

	src, err := file.Open()
	if err != nil {
		return nil, false, errors.New("failed to open file")
	}
	defer src.Close()

	dst, err := os.CreateTemp("", "uploaded-*.png")
	if err != nil {
		return nil, false, errors.New("failed to create temporary file")
	}
	defer func() {
		dst.Close()
		os.Remove(dst.Name())
	}()

	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, false, errors.New("failed to copy file contents")
	}

	inputImage := gocv.IMRead(dst.Name(), gocv.IMReadColor)
	if inputImage.Empty() {
		return nil, false, errors.New("failed to read uploaded image")
	}
	defer inputImage.Close()

	tempFile, err := os.CreateTemp("", "output-*.jpg")
	if err != nil {
		return nil, false, errors.New("failed to create temporary file")
	}

	if ok := gocv.IMWrite(tempFile.Name(), inputImage); !ok {
		return nil, false, errors.New("failed to write output image")
	}
	return tempFile, true, nil
}
