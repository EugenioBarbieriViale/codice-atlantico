package main

// TODO
// - fix float64 json

import (
	"fmt"
	"log"
	"encoding/json"
	"io"
	"net/http"

	"github.com/EugenioBarbieriViale/codice-atlantico/database"
)

func main() {
	cfg := database.DefaultConfig()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// b, err := database.NewBook("test1", "test2", "test3", 0.0, "test4")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// booksContent, err := database.ShowTable[database.Book](db, "books")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("%v\n", booksContent)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})

	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		book := getBook(w, r)
		if book.Title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}
		fmt.Println(book)

		// err := db.AddBook(book)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Fprintf(w, "Book added")
	})

	port := "5000" 
	fmt.Println("server running on port", port)

	err = http.ListenAndServe(":" + port, nil)
	if err != nil {
		log.Fatal(err)
	}
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
