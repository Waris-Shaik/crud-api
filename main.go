package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// rendering .env file
func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error in .env🔒🔑 File")
	}
}

func main() {

	// Setting PORT
	PORT := os.Getenv("PORT")

	// if port is empty throw error anad return .
	if PORT == "" {
		log.Fatal("Error PORT is missing in .env File")
	}

	// router
	router := mux.NewRouter()

	// routes
	router.HandleFunc("/", HomeHandler).Methods("GET")
	router.HandleFunc("/movies/new", CreateMovie).Methods("POST")
	router.HandleFunc("/movies", GetAllMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", GetMoivieById).Methods("GET")
	router.HandleFunc("/movies/{id}", DeleteMovie).Methods("DELETE")
	router.HandleFunc("/movies/{id}", UpdateMovie).Methods("PUT")

	// server listening
	fmt.Printf("Server is listening on PORT: %v\n", PORT)
	if err := http.ListenAndServe(":"+PORT, router); err != nil {
		log.Fatal("Error in listening server: ", err)
	}

}
