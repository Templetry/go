// Command template-app is the TemplateApp CLI.
package main

import (
	"flag"
	"fmt"
	"os"

	"example.com/template-app/internal/greet"
)

// version is stamped by release builds via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	name := flag.String("name", "world", "who to greet")
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Println("template-app", version)
		os.Exit(0)
	}
	fmt.Println(greet.Greeting(*name))
}
