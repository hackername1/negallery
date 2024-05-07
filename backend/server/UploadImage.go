package server

import (
	"fmt"
	"io"
	"log"
	"main/processor"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {
	// Get image parameters
	var image processor.Image
	var imageData ImageData

	// Check password
	if r.FormValue("password") != myEnvironment["GALLERY_ADMIN_PASSWORD"] {
		http.Error(w, "Wrong password", http.StatusUnauthorized)
		return
	}

	var title = r.FormValue("title")
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	imageData.Title = title
	imageData.Country = r.FormValue("country")
	imageData.Date = r.FormValue("date")

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, handler, err := r.FormFile("raw_image")
	if err != nil {
		http.Error(w, serverErrorString, http.StatusBadRequest)
		log.Println("Error getting the file from the form: ", err)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Println("Error closing the file: ", err)
		}
	}(file)

	// Check the file format to match jpg or png
	if handler.Header.Get("Content-Type") != "image/jpeg" &&
		handler.Header.Get("Content-Type") != "image/png" {
		http.Error(w, "Wrong image format", http.StatusBadRequest)
		return
	}

	// Read the file
	image.ImageBytes, err = io.ReadAll(file)

	// Process the image
	imagePaths, successfulWidths, err := processImage(image)
	if err != nil {
		http.Error(w, serverErrorString, http.StatusBadRequest)
		return
	}

	imageData.Resolutions = make([]string, len(successfulWidths))
	for i, width := range successfulWidths {
		imageData.Resolutions[i] = strconv.Itoa(width)
	}

	// Add record to database
	index, err := writeToDatabase(imageData)
	if err != nil {
		http.Error(w, serverErrorString, http.StatusBadRequest)
		log.Println("Error writing to the database: ", err)
		return
	}

	// Make paths for the S3
	var endpointPaths = make([]string, len(imagePaths))
	for i, successfulWidth := range successfulWidths {
		endpointPaths[i] = myEnvironment["S3_BASE_PATH"] + "/" +
			fmt.Sprintf("%06d", index) + "/" +
			strconv.Itoa(successfulWidth) +
			".jpg"
	}

	// Load image to S3
	err = loadToS3(imagePaths, endpointPaths)
	if err != nil {
		http.Error(w, serverErrorString, http.StatusBadRequest)
		return
	}

	// Delete everything in the temp directory
	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Println("Error getting the current directory: ", err)
	} else {
		err = os.RemoveAll(currentDirectory + tempPath)
		if err != nil {
			log.Println("Error removing the temp directory: ", err)
		}
	}

	// Return success
	_, err = w.Write([]byte("Success"))
	if err != nil {
		http.Error(w, serverErrorString, http.StatusBadRequest)
		log.Println("Error writing the response: ", err)
		return
	}
}

func processImage(image processor.Image) (imagePaths []string, successfulWidths []int, err error) {
	var currentUnixTime = time.Now().Unix()
	var currentUnixTimeString = strconv.FormatInt(currentUnixTime, 10)

	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Println("Error getting the current directory: ", err)
		return nil, nil, err
	}

	var pathBase = currentDirectory + tempPath + "/" + currentUnixTimeString + "/"

	// Create directory

	err = os.MkdirAll(pathBase, 0777)
	if err != nil {
		log.Println("Error creating the directory: ", err)
		return nil, nil, err
	}

	for _, width := range widths {
		// Resize the image
		resizedImage, err := image.ResizeByWidth(width)
		if err == processor.ImageIsSmallerError {
			continue
		} else if err != nil {
			log.Println("Error processing the image: ", err)
			return nil, nil, err
		}

		// Save the image
		var imagePath = pathBase
		if width != 0 {
			imagePath += strconv.Itoa(width) + ".jpg"
		} else {
			imagePath += "origin.jpg"
		}

		err = resizedImage.Save(imagePath)
		if err != nil {
			log.Println("Error saving the image: ", err)
			return nil, nil, err
		}

		imagePaths = append(imagePaths, imagePath)
		successfulWidths = append(successfulWidths, width)
	}

	return imagePaths, successfulWidths, nil
}

func writeToDatabase(imageData ImageData) (index int, err error) {
	// Get the biggest index
	err = database.QueryRow("SELECT MAX(`index`) FROM images").Scan(&index)
	if err != nil {
		log.Println("Error getting the biggest index: ", err)
		return 0, err
	}
	index++

	// Create resolutions string
	var resolutionsString = ""
	for i := 0; i < len(imageData.Resolutions); i++ {
		if imageData.Resolutions[i] == "0" {
			resolutionsString += "origin"
		} else {
			resolutionsString += imageData.Resolutions[i]
		}
		if i != len(imageData.Resolutions)-1 {
			resolutionsString += ","
		}
	}

	// Insert the image data
	_, err = database.Exec(
		"INSERT INTO images (`index`, title, country, date, resolutions) VALUES (?, ?, ?, ?, ?)",
		index,
		imageData.Title,
		imageData.Country,
		imageData.Date,
		resolutionsString,
	)
	if err != nil {
		log.Println("Error inserting the image data: ", err)
		return
	}

	return index, nil
}
