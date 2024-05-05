package server

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func Run() {
	// Determine whether host is local or web
	if _, err := os.Stat(myEnvironment["CERT_FILE"]); os.IsNotExist(err) {
		localMode = true
	}

	// Load the environment variables
	loadEnvironment()

	// Connect to the database
	database = ConnectDatabase()

	// Create amd set up the router
	router := mux.NewRouter()
	router.Use(corsMiddleware)
	router.Use(requestCheck)

	// Define the routes
	router.HandleFunc(apiUrl+"/getImage", GetImage).Methods("GET")
	router.HandleFunc(apiUrl+"/getImageList", GetImageList).Methods("GET")
	router.HandleFunc(apiUrl+"/uploadImage", UploadImage).Methods("POST")

	// Start the server
	startServer(router)
}

// Connect to the database
func ConnectDatabase() *sql.DB {
	var db *sql.DB
	configuration := mysql.NewConfig()
	(*configuration).Net = "tcp"
	(*configuration).Addr = myEnvironment["GALLERY_HOST"]
	(*configuration).User = myEnvironment["GALLERY_USER"]
	(*configuration).Passwd = myEnvironment["GALLERY_PASSWORD"]
	(*configuration).DBName = myEnvironment["GALLERY_DATABASE"]
	(*configuration).ParseTime = true

	db, err := sql.Open("mysql", configuration.FormatDSN())
	if err != nil {
		panic(err)
	}

	return db
}

// Load the environment variables
func loadEnvironment() {
	err := godotenv.Load(envPath)
	if err != nil {
		panic(err)
	}
	myEnvironment, err = godotenv.Read()
	if err != nil {
		panic(err)
	}
}

// Define the CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers to allow cross-origin requests
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// If it's a preflight request, send a 200 OK status
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// Define the request check middleware
func requestCheck(next http.Handler) http.Handler {
	// Set up the router mode variables
	var host = webHost
	var scheme = "https"
	if localMode {
		host = localHost
		scheme = "http"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request host
		if r.Host != host {
			http.Error(w, "404 not found", http.StatusNotFound)
			return
		}

		// Check the request scheme
		if r.URL.Scheme != scheme {
			http.Error(w, "Invalid scheme", http.StatusForbidden)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// Start the server
func startServer(router *mux.Router) {
	if localMode {
		// Start the HTTP server
		log.Fatal(
			http.ListenAndServe(
				serverPort,
				router,
			))
	} else {
		// Start the HTTPS server
		log.Fatal(
			http.ListenAndServeTLS(
				serverPort,
				myEnvironment["CERT_FILE"],
				myEnvironment["PRIVATE_KEY_FILE"],
				router,
			))
	}
}
