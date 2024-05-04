package server

import "database/sql"

// Variables
var database *sql.DB
var myEnvironment map[string]string

// Constants
const envPath = ".env"
const tempPath = "/temp"
const serverPort = ":8443"
const apiUrl = "/api/v1"
const serverErrorString = "Error on the server"

const (
	widthOrigin = 0
	width1920   = 1920
	width1280   = 1280
	width640    = 640
)

var widths = []int{
	widthOrigin,
	width1920,
	width1280,
	width640,
}
