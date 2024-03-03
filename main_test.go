package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var jobs = []Job{
	{
		InputPath:  "input1",
		OutPath:    "output1",
		ReadTime:   1.0,
		ResizeTime: 2.0,
		GrayTime:   3.0,
		WriteTime:  4.0,
	},
	{
		InputPath:  "input2",
		OutPath:    "output2",
		ReadTime:   1.1,
		ResizeTime: 2.2,
		GrayTime:   3.3,
		WriteTime:  4.4,
	},
}

func TestJobsToJSON(t *testing.T) {
	filename := "test.json"
	jobsToJSON(jobs, filename)

	bytes, err := os.ReadFile(filename)
	assert.NoError(t, err)

	var existingJobs []Job
	err = json.Unmarshal(bytes, &existingJobs)
	assert.NoError(t, err)

	assert.Equal(t, len(jobs), len(existingJobs))

	os.Remove(filename)
}

func TestLoadProcessedPaths(t *testing.T) {
	filename := "test.json"
	jobsToJSON(jobs, filename)

	processedPaths := loadProcessedPaths(filename)

	assert.Equal(t, len(jobs), len(processedPaths))

	for _, job := range jobs {
		assert.True(t, processedPaths[job.InputPath])
	}

	os.Remove(filename)
}
