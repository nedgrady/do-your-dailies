package main

import (
	"log"
	"net/http"

	"do-your-dailies/server/internal/api"
	"do-your-dailies/server/internal/db"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=dailies port=5432 sslmode=disable"

	database, err := db.New(dsn)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	app := api.New(database)

	log.Println("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", app.Router))
}
