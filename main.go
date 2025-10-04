package main

import (
	"fmt"
	"log"

	"database/sql"
	_ "github.com/lib/pq"

	"github.com/EugenioBarbieriViale/codice-atlantico/config"
)

func main() {
	cfg := config.NewConfig()
	args := cfg.ToString()

	fmt.Println("Connecting to database...")
	db, err := sql.Open("postgres", args)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to database")
}
