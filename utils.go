package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"_id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Hero     string    `json:"hero"`
	Director *Director `json:"director"`
}

type Director struct {
	Name string `json:"name"`
}

type Movies []Movie

var movies = Movies{
	Movie{ID: "1011", Title: "OG", Content: "A Gangster Backdrop Film", Hero: "Pawan Kalyan", Director: &Director{Name: "Sujith"}},
	Movie{ID: "1012", Title: "Game Changer", Content: "A Political Backdrop Film", Hero: "Ram Charan", Director: &Director{Name: "Shankar"}},
}

// Handlers

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Invalid Path", http.StatusNotFound)
		return
	}

	type Response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	data := Response{
		Success: true,
		Message: "Welcome to GoLang Server",
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Println("Error encoding JSON:", err)
		http.Error(w, "Error While Encoding Data", http.StatusInternalServerError)
		return
	}

}

func GetAllMovies(w http.ResponseWriter, r *http.Request) {

	type Response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Movies  Movies `json:"movies"`
	}

	data := Response{
		Success: true,
		Message: "Successfully got data",
		Movies:  movies,
	}

	// set Headers
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Println("Error encoding JSON:", err)
		http.Error(w, "Error While Encoding Data", http.StatusInternalServerError)
		return
	}

}

func CreateMovie(w http.ResponseWriter, r *http.Request) {

	var movie Movie
	// fmt.Println(r.Body)

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Error while getting data from client", 400)
		return
	}

	if movie.ID == "" || movie.Title == "" || movie.Content == "" || movie.Hero == "" || movie.Director == nil || movie.Director.Name == "" {
		errResponse := struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: false,
			Message: "Please fill all required field's",
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	// fmt.Println(movie.ID, "here u go bro id is")

	// check if id is already exists
	for _, item := range movies {
		if item.ID == movie.ID {
			errResponse := struct {
				Success bool   `json:"success"`
				Message string `json:"message"`
			}{
				Success: false,
				Message: "Id Already Exists it should be unique",
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errResponse)
			return
		}
	}

	// fmt.Println(movie)

	// appending in  movies
	movies = append(movies, movie)
	successresponse := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Movie   Movie  `json:"movie"`
	}{
		Success: true,
		Message: "Successfully movie created",
		Movie:   movie,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(successresponse)
}

func GetMoivieById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Id should not be null", 400)
		return
	}

	var movieToBeShown *Movie

	// check if ID exsits or not
	for i, item := range movies {
		if item.ID == id {
			movieToBeShown = &movies[i]
			break
		}
	}

	if movieToBeShown == nil {
		response := struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: false,
			Message: "Inavlid Id Movie not foud",
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Movie   Movie  `json:"movie"`
	}{
		Success: true,
		Message: "Succesfully found movie",
		Movie:   *movieToBeShown,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return

}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]

	if !ok {
		http.Error(w, "Id Should not be null", 400)
		return
	}

	// check if ID exists
	var movieToDelete *Movie

	for i, item := range movies {
		if item.ID == id {
			movieToDelete = &movies[i]
		}
	}

	if movieToDelete == nil {
		response := struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: false,
			Message: "Invalid Id Movie not found",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	var showDeletedMovie Movie

	for i, item := range movies {
		if item.ID == id {
			showDeletedMovie = item
			// deleting movie
			movies = append(movies[:i], movies[i+1:]...)
			break
		}
	}

	response := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Movie   Movie  `json:"movie"`
	}{
		Success: true,
		Message: "Successfully Movie Deleted",
		Movie:   showDeletedMovie,
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		http.Error(w, "Id should not be null", http.StatusBadRequest)
		return
	}

	// check if id exists or not if not throw error and return handler
	var movieToUpdate *Movie

	for i, item := range movies {
		if item.ID == id {
			movieToUpdate = &movies[i]
			break
		}
	}

	// fmt.Println(movieToUpdate)

	if movieToUpdate == nil {
		response := struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: false,
			Message: "Invalid Id Movie not found",
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)

		return
	}

	var updateMovie Movie

	if err := json.NewDecoder(r.Body).Decode(&movieToUpdate); err != nil {
		fmt.Println("Error parsing JSON:", err)
		http.Error(w, "Error Parsing JSON", 400)
		return
	}

	if updateMovie.Title != "" {
		movieToUpdate.Title = updateMovie.Title
	}
	if updateMovie.Content != "" {
		movieToUpdate.Content = updateMovie.Content
	}

	if updateMovie.Hero != "" {
		movieToUpdate.Hero = updateMovie.Hero
	}

	if updateMovie.Director != nil && updateMovie.Director.Name != "" {
		movieToUpdate.Director = &Director{Name: updateMovie.Director.Name}
	}

	successresponse := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Movie   Movie  `json:"movie"`
	}{
		Success: true,
		Message: "Successfully Movie Updated",
		Movie:   *movieToUpdate,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successresponse)

}
