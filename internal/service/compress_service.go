package service

import (
	"context"
	"errors"
	"go-img-processing/bootstrap"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"gocv.io/x/gocv"
)

type CompressService interface {
	Compress(ctx context.Context, src *multipart.FileHeader, quality int) (*os.File, bool, error)
}

type CompressServiceImpl struct {
	cfg *bootstrap.Container
}

func NewCompressService(cfg *bootstrap.Container) CompressService {
	return &CompressServiceImpl{cfg: cfg}
}

func (s *CompressServiceImpl) Compress(ctx context.Context, file *multipart.FileHeader, quality int) (*os.File, bool, error) {
	extArray := []string{".png", ".jpg", ".jpeg"}
	extension := filepath.Ext(file.Filename)
	extList := strings.Join(extArray, " ")
	if found := strings.Contains(extList, extension); !found {
		return nil, false, errors.New("only PNG, JPG, JPEG files are accepted")
	}

	src, err := file.Open()
	if err != nil {
		return nil, false, errors.New("failed to open file")
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		return nil, false, errors.New("failed to decode image: " + err.Error())
	}

	mat, err := gocv.ImageToMatRGB(img)
	if err != nil {
		return nil, false, errors.New("failed to convert image to Mat: " + err.Error())
	}
	defer mat.Close()

	buffer, err := gocv.IMEncodeWithParams(gocv.FileExt(extension), mat, []int{gocv.IMWriteJpegQuality, quality})
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
