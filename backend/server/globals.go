package server

import "database/sql"

// Variables
var database *sql.DB
var myEnvironment map[string]string

// Constants
const envPath = ".env"
const serverPort = ":8443"
const apiUrl = "/api/v1"
