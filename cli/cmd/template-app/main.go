// Command template-app is the TemplateApp CLI.
package main

import (
	"flag"
	"fmt"
	"os"

	// tpl:if environments
	"example.com/template-app/internal/config"
	// tpl:endif
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

	// tpl:if environments
	// No `env` subcommand: adding one would change the CLI's shape, which a
	// feature must not do. The profile changes what the existing run does —
	// diagnostics go to stderr, so piping stdout stays clean.
	if cfg, err := config.Load(".", ""); err == nil && cfg.VerboseErrors {
		fmt.Fprintf(os.Stderr, "[%s] log level %s\n", cfg.Environment, cfg.LogLevel)
	}
	// tpl:endif

	fmt.Println(greet.Greeting(*name))
}
