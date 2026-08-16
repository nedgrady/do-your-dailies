package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"do-your-dailies/server/internal/api"
	"do-your-dailies/server/internal/db"
	"do-your-dailies/server/internal/migrations"
)

func resolveDSN(getenv func(string) string) (string, error) {
	dsn := getenv("DATABASE_URL")
	if dsn == "" {
		return "", errors.New("DATABASE_URL is not set")
	}
	return dsn, nil
}

func resolvePort(getenv func(string) string) string {
	if port := getenv("PORT"); port != "" {
		return port
	}
	return "8080"
}

func main() {
	time.Local = time.UTC

	dsn, err := resolveDSN(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.New(dsn)

	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	if err := migrations.Migrate(database); err != nil {
		log.Fatal("failed to run migrations:", err)
	}

	app := api.New(database)

	port := resolvePort(os.Getenv)
	log.Printf("starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, app.Router))
}
