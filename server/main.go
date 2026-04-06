package main

import (
	"do-your-dailies/server/api"
	"log"
	"net/http"
)

func main() {
	app := api.New()

	log.Println("starting server on :8080")
	err := http.ListenAndServe(":8080", app.Router)
	if err != nil {
		log.Fatal(err)
	}
}
