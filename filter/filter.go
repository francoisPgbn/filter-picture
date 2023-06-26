package filter

import (
	"github.com/disintegration/imaging"
)

type Filter interface {
	Process(srcPath, dstPath string) error
}

type Grayscale struct {
}

type Blur struct {
}

func (g Grayscale) Process(srcPath, destPath string) error {
	srcImg, err := imaging.Open(srcPath)
	if err != nil {
		return err
	}

	grayscaleImg := imaging.Grayscale(srcImg)

	err = imaging.Save(grayscaleImg, destPath)
	if err != nil {
		return err
	}

	return nil
}

func (g Blur) Process(srcPath, destPath string) error {
	srcImg, err := imaging.Open(srcPath)
	if err != nil {
		return err
	}

	grayscaleImg := imaging.Blur(srcImg, 4)

	err = imaging.Save(grayscaleImg, destPath)
	if err != nil {
		return err
	}

	return nil
}
