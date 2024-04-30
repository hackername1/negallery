package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetImageData(writer http.ResponseWriter, request *http.Request) {
	var response Response

	// Get arguments from URL
	queryParams := request.URL.Query()
	imageId := queryParams.Get("id")

	imageIdInt, err := strconv.Atoi(imageId)
	if err != nil {
		response = Response{
			Status:   "error",
			Error:    "invalid image ID",
			Response: "",
		}
		json.NewEncoder(writer).Encode(response)
	}

	imageData, err := dbGetImageData(imageIdInt)
	response = Response{
		Status:   "success",
		Error:    err.Error(),
		Response: imageData,
	}

	json.NewEncoder(writer).Encode(response)
}

func dbGetImageData(imageId int) (ImageData, error) {
	var query = "SELECT * FROM images WHERE `index` = ?"

	// Execute the query
	var rows, err = database.Query(query, imageId)
	if err != nil {
		panic(err)
	}

	// Get the data
	var imageData ImageData
	var resolutionsString string
	for rows.Next() {
		if err = rows.Scan(&imageData.Id, &imageData.Title, &imageData.Country, &imageData.Date, &resolutionsString); err != nil {
			log.Println(err)
			return ImageData{}, errors.New("image not found")
		}
		resolutions := strings.Split(resolutionsString, ",")
		imageData.Resolutions = resolutions
	}

	if imageData.Resolutions == nil {
		return ImageData{}, errors.New("image not found")
	}

	return imageData, errors.New("")
}
