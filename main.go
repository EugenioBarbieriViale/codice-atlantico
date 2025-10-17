package main

import (
	"fmt"
	"log"

	"github.com/EugenioBarbieriViale/codice-atlantico/database"
	"github.com/EugenioBarbieriViale/codice-atlantico/server"
)

func main() {
	cfg := database.DefaultConfig()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	port := "5000"
	server.StartHTTPServer(db, port)

	fmt.Println("Done")
}
