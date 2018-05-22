package main

import (
	"image/jpeg"
	"io"

	"github.com/nfnt/resize"
)

func Resize(src io.Reader, width, height int, dest io.Writer) error {
	img, err := jpeg.Decode(src)
	if err != nil {
		return err
	}

	resizedImg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	return jpeg.Encode(dest, resizedImg, nil)
}
