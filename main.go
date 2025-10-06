package main

import (
	"fmt"
	"log"
	// "net/http"

	"github.com/EugenioBarbieriViale/codice-atlantico/database"
)

func main() {
	cfg := database.DefaultConfig()

	db, err := database.Connect(cfg)
	check(err)
	defer db.Close()

	b, err := database.NewBook("test1", "test2", "test3", 0.0, "test4")
	check(err)

	err = db.AddBook(b)
	check(err)

	booksContent, err := database.ShowTable[database.Book](db, "books")
	check(err)

	fmt.Println(booksContent)

	// http.Handle("/", http.FileServer(http.Dir("./static")))

	// http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Hello")
	// })

	// port := "5000" 
	// fmt.Println("server running on port", port)

	// log.Fatal(http.ListenAndServe(":" + port, nil))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

