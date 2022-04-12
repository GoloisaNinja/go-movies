package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func sortMovies() []Movie {
	sort.Slice(
		movies, func(i, j int) bool {
			return movies[i].ID < movies[j].ID
		},
	)
	return movies
}

func movieIdExists(id string) bool {
	for _, movie := range movies {
		if movie.ID == id {
			return true
		}
	}
	return false
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sortMovies())
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	if !movieIdExists(params["id"]) {
		http.Error(w, "404 movie does not exist...", http.StatusNotFound)
	}
	for _, movie := range movies {
		if movie.ID == params["id"] {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	if !movieIdExists(params["id"]) {
		http.Error(w, "404 movie does not exist...", http.StatusNotFound)
		return
	}
	for i, movie := range movies {
		if movie.ID == params["id"] {
			movies = append(movies[:i], movies[i+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(100000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	if !movieIdExists(params["id"]) {
		http.Error(w, "404 no movie found...", http.StatusNotFound)
		return
	}
	var updatedMovie Movie
	_ = json.NewDecoder(r.Body).Decode(&updatedMovie)
	updatedMovie.ID = params["id"]

	for i, movie := range movies {
		if movie.ID == params["id"] {
			movies = append(movies[:i], movies[i+1:]...)
			movies = append(movies, updatedMovie)
			json.NewEncoder(w).Encode(updatedMovie)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()

	movies = append(
		movies, Movie{
			ID:       "1",
			Isbn:     "222111",
			Title:    "Movie One",
			Director: &Director{Firstname: "John", Lastname: "Doe"},
		},
	)
	movies = append(
		movies, Movie{
			ID:       "2",
			Isbn:     "111222",
			Title:    "Movie Two",
			Director: &Director{Firstname: "LouLou", Lastname: "Collins"},
		},
	)
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
