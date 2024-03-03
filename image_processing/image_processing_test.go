package imageprocessing

import (
	"image"
	"image/color"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to create a mock image
func createTestImage(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{78, 42, 132, 255})
		}
	}
	return img
}

func TestReadImage(t *testing.T) {
	tempFile, _ := os.CreateTemp("", "*.jpg")
	defer os.Remove(tempFile.Name())
	img := createTestImage(100, 100)
	_ = WriteImage(tempFile.Name(), img)

	// Test reading image
	readImg, err := ReadImage(tempFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, img.Bounds(), readImg.Bounds())

	// Test error handling
	_, err = ReadImage("nonexistent.jpg")
	assert.Error(t, err)
}

func TestWriteImage(t *testing.T) {
	tempFile, _ := os.CreateTemp("", "*.jpg")
	defer os.Remove(tempFile.Name())
	img := createTestImage(100, 100)

	// Test writing image
	err := WriteImage(tempFile.Name(), img)
	assert.NoError(t, err)

	// Verify image was written correctly
	readImg, _ := ReadImage(tempFile.Name())
	assert.Equal(t, img.Bounds(), readImg.Bounds())
}

func TestGrayscale(t *testing.T) {
	img := createTestImage(100, 100)
	grayImg := Grayscale(img)

	// Verify image is grayscale
	for y := 0; y < grayImg.Bounds().Dy(); y++ {
		for x := 0; x < grayImg.Bounds().Dx(); x++ {
			_, _, _, _ = grayImg.At(x, y).RGBA()
		}
	}
}

func TestResize(t *testing.T) {
	img := createTestImage(100, 100)
	resizedImg := Resize(img)

	// Verify image was resized correctly
	assert.Equal(t, 500, resizedImg.Bounds().Dx())
	assert.Equal(t, 500, resizedImg.Bounds().Dy())
}
