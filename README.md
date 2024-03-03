# Image Processing in Go

Submitted assignment for Northwestern MSDS 431

## Overview

This repository is a fork of an image processing pipeline set up by [code-heim](https://github.com/code-heim/go_21_goroutines_pipeline), making the following changes:

- Error handling for image inputs and outputs, specifically checking for file-not-found and permission errors. 
- Pictures of cats to test the image processing pipeline
- Added unit tests for critical components
- Added benchmarking for each step in the pipeline, adding time features to the Job struct
- Writing results to JSON, checking for unique filepaths before appending.
- Added rotation and blur to the processing pipeline using the [imaging package](https://github.com/disintegration/imaging)


## Usage

- Clone the repository 
- In the repository directory run `go build .`
- Add any new images to the `images` directory
- run the command `.\goroutines_pipeline.exe`

This program will store processing times for all new images in the `images` directory of the repository, in other words an image will not be processed if it already has a record in the `results.json` file that gets generated after the first run of the program. New images will still be processed even if duplicates exist in the `images` directory. 

## Photo Credits

- `cat_1.jpg` -- "Cat" by London's is licensed under CC BY-SA 2.0.
- `cat_2.jpg` -- "cats group photo" by S@veOurSm:)e is licensed under CC BY-SA 2.0.
- `cat_3.jpg` -- "Cats sleep anywhere." by rawdonfox is licensed under CC BY 2.0.
- `cat_4.jpg` -- "cat" by wapiko☆ is licensed under CC BY 2.0.