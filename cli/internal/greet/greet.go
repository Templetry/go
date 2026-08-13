// Package greet is the sample library package the tests exercise.
package greet

import "fmt"

// Greeting builds the message the CLI prints.
func Greeting(name string) string {
	if name == "" {
		name = "world"
	}
	return fmt.Sprintf("Hello, %s!", name)
}
