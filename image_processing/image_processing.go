package imageprocessing

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
)

func ReadImage(path string) (image.Image, error) {
	inputFile, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("File not found: %s", path)
		}
		return nil, fmt.Errorf("failed to open image file %s: %w", path, err)
	}
	defer inputFile.Close()

	// Decode the image
	img, _, err := image.Decode(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image %s: %w", path, err)
	}
	return img, nil
}

func WriteImage(path string, img image.Image) error {
	outputFile, err := os.Create(path)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			log.Printf("Permission denied: %s", path)
		}
		return fmt.Errorf("failed to create output file %s: %w", path, err)
	}
	defer outputFile.Close()

	// Encode the image to the new file
	err = jpeg.Encode(outputFile, img, nil)
	if err != nil {
		return fmt.Errorf("failed to encode image %s: %w", path, err)
	}
	return nil
}

func Grayscale(img image.Image) image.Image {
	// Create a new grayscale image
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	// Convert each pixel to grayscale
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalPixel := img.At(x, y)
			grayPixel := color.GrayModel.Convert(originalPixel)
			grayImg.Set(x, y, grayPixel)
		}
	}
	return grayImg
}

func Resize(img image.Image) image.Image {
	newWidth := uint(500)
	newHeight := uint(500)
	resizedImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
	return resizedImg
}

func Rotate(img image.Image, angle float64) image.Image {
	rotatedImg := imaging.Rotate(img, angle, color.Transparent)
	return rotatedImg
}

func Blur(img image.Image, sigma float64) image.Image {
	blurredImg := imaging.Blur(img, sigma)
	return blurredImg
}
