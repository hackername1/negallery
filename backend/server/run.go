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
	loadEnvironment()
	database = ConnectDatabase()
	router := mux.NewRouter()
	router.Use(corsMiddleware)

	router.HandleFunc(apiUrl+"/getImage", GetImage).Methods("GET")
	router.HandleFunc(apiUrl+"/getImageList", GetImageList).Methods("GET")
	router.HandleFunc(apiUrl+"/uploadImage", UploadImage).Methods("POST")

	startServer(router)
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

// Start the server
func startServer(router *mux.Router) {
	var sqlMode = false
	if _, err := os.Stat("fullchain.crt"); os.IsNotExist(err) {
		sqlMode = true
	}

	if sqlMode {
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
				myEnvironment["KEY_FILE"],
				router,
			))
	}
}
