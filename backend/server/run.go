package server

import (
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

	router.HandleFunc(apiUrl+"/getImageData", GetImageData).Methods("GET")

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
		// Start the HTTP server
		log.Fatal(
			http.ListenAndServeTLS(
				serverPort,
				myEnvironment["CERT_FILE"],
				myEnvironment["KEY_FILE"],
				router,
			))
	}
}