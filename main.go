package main

import (
	"fmt"
	imageprocessing "goroutines_pipeline/image_processing"
	"image"
	"strings"
)

type Job struct {
	InputPath string
	Image     image.Image
	OutPath   string
	Error     error
}

func loadImage(paths []string) <-chan Job {
	out := make(chan Job)
	go func() {
		// For each input path create a job and add it to
		// the out channel
		for _, p := range paths {
			job := Job{InputPath: p,
				OutPath: strings.Replace(p, "images/", "images/output/", 1)}
			job.Image, job.Error = imageprocessing.ReadImage(p)
			if job.Error != nil {
				fmt.Println("Error reading image:", job.Error)
				continue
			}
			out <- job
		}
		close(out)
	}()
	return out
}

func resize(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		// For each input job, create a new job after resize and add it to
		// the out channel
		for job := range input { // Read from the channel
			job.Image = imageprocessing.Resize(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

func convertToGrayscale(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input { // Read from the channel
			job.Image = imageprocessing.Grayscale(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

func saveImage(input <-chan Job) <-chan bool {
	out := make(chan bool)
	go func() {
		for job := range input { // Read from the channel
			err := imageprocessing.WriteImage(job.OutPath, job.Image)
			if err != nil {
				fmt.Println("Error writing image:", err)
				out <- false
				continue
			}
			out <- true
		}
		close(out)
	}()
	return out
}

func main() {

	imagePaths := []string{"images/cat_1.jpg",
		"images/cat_2.jpg",
		"images/cat_3.jpg",
		"images/cat_4.jpg",
	}

	channel1 := loadImage(imagePaths)
	channel2 := resize(channel1)
	channel3 := convertToGrayscale(channel2)
	writeResults := saveImage(channel3)

	for success := range writeResults {
		if success {
			fmt.Println("Success!")
		} else {
			fmt.Println("Failed!")
		}
	}
}
