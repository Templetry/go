// Command server runs the TemplateApp HTTP service.
package main

import (
	"log"
	"net/http"
	"os"

	"example.com/template-app/internal/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("template-app listening on :%s", port)
	if err := http.ListenAndServe(":"+port, api.NewMux()); err != nil {
		log.Fatal(err)
	}
}
