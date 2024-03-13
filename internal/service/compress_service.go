package service

import (
	"context"
	"errors"
	"go-img-processing/bootstrap"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"gocv.io/x/gocv"
)

type CompressService interface {
	Compress(ctx context.Context, src io.Reader, quality int) (*os.File, bool, error)
}

type CompressServiceImpl struct {
	cfg *bootstrap.Container
}

func NewCompressService(cfg *bootstrap.Container) CompressService {
	return &CompressServiceImpl{cfg: cfg}
}

func (s *CompressServiceImpl) Compress(ctx context.Context, file io.Reader, quality int) (*os.File, bool, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, false, errors.New("failed to decode image: " + err.Error())
	}

	mat, err := gocv.ImageToMatRGB(img)
	if err != nil {
		return nil, false, errors.New("failed to convert image to Mat: " + err.Error())
	}
	defer mat.Close()

	buffer, err := gocv.IMEncodeWithParams(gocv.JPEGFileExt, mat, []int{gocv.IMWriteJpegQuality, quality})
	if err != nil {
		return nil, false, errors.New("failed to encode image: " + err.Error())
	}
	defer buffer.Close()

	tempFile, err := os.CreateTemp("", "compressed-*.jpg")
	if err != nil {
		return nil, false, errors.New("failed to create temporary file")
	}
	defer tempFile.Close()

	_, err = tempFile.Write(buffer.GetBytes())
	if err != nil {
		return nil, false, errors.New("failed to write output image")
	}

	return tempFile, true, nil
}
