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

type ResizeService interface {
	Resize(ctx context.Context, file io.Reader, width int, height int, ar bool, scalar *gocv.Scalar, extension string) (*os.File, bool, error)
}

type ResizeServiceImpl struct {
	cfg *bootstrap.Container
}

func NewResizeService(cfg *bootstrap.Container) ResizeService {
	return &ResizeServiceImpl{cfg: cfg}
}

func (s *ResizeServiceImpl) Resize(ctx context.Context, file io.Reader, width int, height int, ar bool, scalar *gocv.Scalar, extension string) (*os.File, bool, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, false, errors.New("failed to decode image: " + err.Error())
	}

	mat, err := gocv.ImageToMatRGB(img)
	if err != nil {
		return nil, false, errors.New("failed to convert image to Mat: " + err.Error())
	}
	defer mat.Close()

	newWidth := width
	newHeight := height
	if ar {
		aspectRatio := float64(mat.Cols()) / float64(mat.Rows())
		if float64(width)/float64(height) > aspectRatio {
			newWidth = int(float64(height) * aspectRatio)
		} else {
			newHeight = int(float64(width) / aspectRatio)
		}
	}

	resized := gocv.NewMat()
	gocv.Resize(mat, &resized, image.Point{X: newWidth, Y: newHeight}, 0, 0, gocv.InterpolationArea)
	defer resized.Close()

	bg := resized
	if scalar != nil {
		bg = gocv.NewMatWithSize(height, width, gocv.MatTypeCV8UC3)
		bg.SetTo(*scalar)
		offsetX := (width - newWidth) / 2
		offsetY := (height - newHeight) / 2
		region := bg.Region(image.Rect(offsetX, offsetY, offsetX+newWidth, offsetY+newHeight))
		resized.CopyTo(&region)
	}

	tempFile, err := os.CreateTemp("", "resized-*"+extension)
	if err != nil {
		return nil, false, errors.New("failed to create temporary file")
	}
	defer tempFile.Close()

	if ok := gocv.IMWrite(tempFile.Name(), bg); !ok {
		return nil, false, errors.New("failed to write output image")
	}

	return tempFile, true, nil
}
