package main

import (
	"encoding/json"
	"fmt"
	imageprocessing "goroutines_pipeline/image_processing"
	"image"
	"log"
	"os"
	"strings"
	"time"
)

const (
	ImagesDir   = "images/"
	ResultsFile = "results.json"
)

type Job struct {
	InputPath  string
	Image      image.Image `json:"-"`
	OutPath    string
	Error      error
	ReadTime   float64
	ResizeTime float64
	GrayTime   float64
	WriteTime  float64
}

// Writes data gathered as part of the Job struct to JSON without the raw image data
func jobsToJSON(jobs []Job, filename string) {
	// Read existing jobs from file
	file, err := os.Open(filename)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error opening file:", err)
		return
	}

	var existingJobs []Job
	if err == nil {
		err = json.NewDecoder(file).Decode(&existingJobs)
		file.Close()
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}
	}

	existingJobs = append(existingJobs, jobs...)

	// Write all jobs back to file
	bytes, err := json.Marshal(existingJobs)
	if err != nil {
		fmt.Println("Error converting jobs to JSON:", err)
		return
	}

	file, err = os.Create(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

// Helper function to prevent duplicate entries in the JSON file.
func loadProcessedPaths(filename string) map[string]bool {
	processedPaths := make(map[string]bool)

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return processedPaths
		}
		fmt.Println("Error opening file:", err)
		return processedPaths
	}
	defer file.Close()

	var jobs []Job
	err = json.NewDecoder(file).Decode(&jobs)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return processedPaths
	}

	for _, job := range jobs {
		processedPaths[job.InputPath] = true
	}

	return processedPaths
}

func loadImage(paths []string, processedPaths map[string]bool) <-chan Job {
	out := make(chan Job)
	go func() {
		// For each input path create a job and add it to
		// the out channel
		for _, p := range paths {
			if _, exists := processedPaths[p]; exists {
				fmt.Printf("Filepath %s already exists\n", p)
				continue
			}
			start := time.Now()
			job := Job{InputPath: p,
				OutPath: strings.Replace(p, "images/", "images/output/", 1)}
			job.Image, job.Error = imageprocessing.ReadImage(p)
			if job.Error != nil {
				fmt.Println("Error reading image:", job.Error)
				continue
			}
			job.ReadTime = time.Since(start).Seconds()
			out <- job
			processedPaths[p] = true
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
			start := time.Now()
			job.Image = imageprocessing.Resize(job.Image)
			job.ResizeTime = time.Since(start).Seconds()
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
			start := time.Now()
			job.Image = imageprocessing.Grayscale(job.Image)
			job.GrayTime = time.Since(start).Seconds()
			out <- job
		}
		close(out)
	}()
	return out
}

func saveImage(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input { // Read from the channel
			start := time.Now()
			err := imageprocessing.WriteImage(job.OutPath, job.Image)
			if err != nil {
				fmt.Println("Error writing image:", err)
				job.Error = err

			} else {
				job.WriteTime = time.Since(start).Seconds()
			}
			out <- job
		}
		close(out)
	}()
	return out
}

func main() {

	dir := ImagesDir
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var imagePaths []string
	for _, file := range files {
		if !file.IsDir() {
			imagePaths = append(imagePaths, dir+file.Name())
		}
	}
	processedPaths := loadProcessedPaths("results.json")
	channel1 := loadImage(imagePaths, processedPaths)
	channel2 := resize(channel1)
	channel3 := convertToGrayscale(channel2)
	writeResults := saveImage(channel3)

	var jobs []Job
	for job := range writeResults {
		if job.Error == nil {
			fmt.Println("Success!")
			jobs = append(jobs, job)
		} else {
			fmt.Println("Failed!")
		}
	}
	jobsToJSON(jobs, "results.json")
}
