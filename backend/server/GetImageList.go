package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetImageList(writer http.ResponseWriter, request *http.Request) {
	var response Response

	// Get arguments from URL
	var queryParams = request.URL.Query()
	var startId = queryParams.Get("start_id")
	var length = queryParams.Get("length")

	// Convert startId to int
	startIdInt, err := strconv.Atoi(startId)
	if err != nil {
		response = Response{
			Status:   "error",
			Error:    "invalid start_id",
			Response: "",
		}
		json.NewEncoder(writer).Encode(response)
		return
	}

	// Convert length to int
	lengthInt, err := strconv.Atoi(length)
	if err != nil {
		response = Response{
			Status:   "error",
			Error:    "invalid length",
			Response: "",
		}
		json.NewEncoder(writer).Encode(response)
		return
	}

	// Get the image data list
	imageDataList, err := GetImageDataListFromDatabase(startIdInt, lengthInt)
	if err != nil {
		response = Response{
			Status:   "error",
			Error:    err.Error(),
			Response: "",
		}
		json.NewEncoder(writer).Encode(response)
		return
	}

	// Send the response
	response = Response{
		Status:   "success",
		Error:    "",
		Response: imageDataList,
	}

	json.NewEncoder(writer).Encode(response)
}

func GetImageDataListFromDatabase(startId int, length int) ([]ImageData, error) {
	var imageDataList []ImageData

	rows, err := database.Query("SELECT * FROM images WHERE `index` >= ? LIMIT ?", startId, length)
	if err != nil {
		return imageDataList, err
	}

	// Iterate over the rows
	for rows.Next() {
		var imageData ImageData
		var resolutionsString string
		if err = rows.Scan(&imageData.Id, &imageData.Title, &imageData.Country, &imageData.Date, &resolutionsString); err != nil {
			log.Println(err)
			return []ImageData{}, errors.New("image not found")
		}
		resolutions := strings.Split(resolutionsString, ",")
		imageData.Resolutions = resolutions

		imageDataList = append(imageDataList, imageData)
	}

	if len(imageDataList) == 0 {
		return []ImageData{}, errors.New("images not found")
	}

	return imageDataList, nil
}
