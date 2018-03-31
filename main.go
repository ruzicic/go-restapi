package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Book Struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Date   string  `json:"date"`
	Title  string  `json:"title"`
	Link   string  `json:"link"`
	Author *Author `json:"author"`
}

// Author Struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Init books var as a slice Book struct
var books []Book

// Get All Books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Get Book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // Get params

	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	json.NewEncoder(w).Encode(&Book{})
}

// Create Book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(1000000))
	books = append(books, book)

	json.NewEncoder(w).Encode(book)
}

// Update Book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			// Delete old book
			books = append(books[:index], books[index+1:]...)

			// Create updated book, with same ID
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]

			// Add new book to slice
			books = append(books, book)

			json.NewEncoder(w).Encode(book)
			return
		}
	}

	json.NewEncoder(w).Encode(books)
}

// Delete Book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(books)
}

func main() {
	// Init Router
	r := mux.NewRouter()

	// Mock Data
	books = append(books, Book{
		ID:    "1",
		Date:  "21 Oct 2014",
		Title: "The Little Go Book",
		Link:  "http://github.com/karlseguin/the-little-go-book",
		Author: &Author{
			Firstname: "Karl",
			Lastname:  "Seguin"}})
	books = append(books, Book{
		ID:    "2",
		Date:  "2012",
		Title: "An Introduction to Programming in Go",
		Link:  "http://www.golang-book.com/",
		Author: &Author{
			Firstname: "Caleb",
			Lastname:  "Doxsey"}})

	// Router Handlers / Endpoints
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
