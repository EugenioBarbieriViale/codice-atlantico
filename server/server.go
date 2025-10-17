package server

import (
	"fmt"
	"log"
	"encoding/json"
	"io"
	"net/http"

	"github.com/EugenioBarbieriViale/codice-atlantico/database"
)

func StartHTTPServer(db *database.Connection, port string) {
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})

	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		handleBook(w, r, db)
	})

	fmt.Println("server running on port", port)

	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleBook(w http.ResponseWriter, r *http.Request, db *database.Connection) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Fatal("Cannot handle POST requests")
	}
	
	book := getBook(w, r)
	if book.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		log.Fatal("Title cannot be empty")
	}
	book.Owner = "test"

	err := db.AddBook(&book)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Book added with ID: %v\n", book.Id)

	booksDb, err := database.GetRow[database.Book](db, book.Id, "books")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", booksDb)
}

func getBook(w http.ResponseWriter, r *http.Request) database.Book {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return database.Book{}
	}
	defer r.Body.Close()
	
	if len(body) == 0 {
		log.Printf("Empty request body")
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return database.Book{}
	}
	
	var book database.Book
	err = json.Unmarshal(body, &book)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return database.Book{}
	}
	
	return book
}
