// Command server runs the TemplateApp REST API.
package main

import (
	"log"
	"net/http"
	"os"

	"example.com/template-app/internal/api"
	"example.com/template-app/internal/store"
)

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "template-app.db"
	}
	db, err := store.Open(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := store.Migrate(db); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("template-app listening on :%s (db %s)", port, dsn)
	if err := http.ListenAndServe(":"+port, api.NewMux(db)); err != nil {
		log.Fatal(err)
	}
}
