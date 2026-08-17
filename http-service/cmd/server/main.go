// Command server runs the TemplateApp HTTP service.
package main

import (
	"log"
	"net/http"
	"os"

	"example.com/template-app/internal/api"
	// tpl:if environments
	"example.com/template-app/internal/config"
	// tpl:endif
)

func main() {
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

	handler := api.NewMux(api.WithEnvironment(cfg.Environment))
	// tpl:endif
	// tpl:if !environments
	handler := api.NewMux()
	// tpl:endif

	log.Printf("template-app listening on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
