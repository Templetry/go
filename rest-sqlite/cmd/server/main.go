// Command server runs the TemplateApp REST API.
package main

import (
	"log"
	"net/http"
	"os"

	"example.com/template-app/internal/api"
	// tpl:if environments
	"example.com/template-app/internal/config"
	// tpl:endif
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

	// tpl:if environments
	// A broken profile stops the service here, before it accepts a request.
	cfg, err := config.Load(".", "")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("template-app starting in %s (log level %s)", cfg.Environment, cfg.LogLevel)

	handler := api.NewMux(db, api.WithEnvironment(cfg.Environment))
	// tpl:endif
	// tpl:if !environments
	handler := api.NewMux(db)
	// tpl:endif

	log.Printf("template-app listening on :%s (db %s)", port, dsn)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
