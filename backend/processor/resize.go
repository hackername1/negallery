package processor

import (
	"bytes"
	"github.com/disintegration/imaging"
	"image"
	"log"
)

func (img *Image) ResizeByWidth(width int) (*Image, error) {
	// Keep the original size
	if width == 0 {
		return &Image{ImageBytes: img.ImageBytes}, nil
	}

	// Create reader for the image bytes
	reader := bytes.NewReader(img.ImageBytes)

	// Decode the image
	image, _, err := image.Decode(reader)
	if err != nil {
		log.Println("Error decoding image: ", err)
		return nil, err
	}

	// Get the image size
	size := image.Bounds().Size()

	scaleFactor := float64(width) / float64(size.X)
	if width > size.X {
		return nil, ImageIsSmallerError
	}

	// Calculate the new height
	var height = int(float64(size.Y) * scaleFactor)

	// Resize the image
	newImage := imaging.Resize(image, width, height, imaging.Lanczos)
	if err != nil {
		return nil, err
	}

	// Get image bytes
	newImageBytes := new(bytes.Buffer)
	err = imaging.Encode(newImageBytes, newImage, imaging.JPEG)
	if err != nil {
		return nil, err
	}

	return &Image{ImageBytes: newImageBytes.Bytes()}, nil
}
