package service

import (
	"context"
	"errors"
	"go-img-processing/bootstrap"
	"io"
	"os"

	"gocv.io/x/gocv"
)

type ConvertService interface {
	Convert(ctx context.Context, file io.Reader) (*os.File, bool, error)
}

type ConvertServiceImpl struct {
	cfg *bootstrap.Container
}

func NewConvertService(cfg *bootstrap.Container) ConvertService {
	return &ConvertServiceImpl{cfg: cfg}
}

func (s *ConvertServiceImpl) Convert(ctx context.Context, file io.Reader) (*os.File, bool, error) {
	dst, err := os.CreateTemp("", "uploaded-*.png")
	if err != nil {
		return nil, false, errors.New("failed to create temporary file")
	}
	defer func() {
		dst.Close()
		os.Remove(dst.Name())
	}()

	_, err = io.Copy(dst, file)
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
