// Package config loads the active environment profile.
//
// Go has no blessed mechanism for this, so the convention is: one
// `.env.<profile>` at the module root, selected by APP_ENV, parsed into a
// validated struct. Nothing else in the service calls os.Getenv.
package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Environments is the closed set of profiles (ADR-0018).
var Environments = []string{"development", "staging", "production"}

// Config is what the service reads.
type Config struct {
	Environment   string
	LogLevel      string
	VerboseErrors bool
	CacheSeconds  int
}

// Load reads a profile by name, or the one APP_ENV selects.
//
// `.env.local` layers on top when present and is gitignored; a real
// environment variable beats both, which is what lets a container run with
// no profile file in the image.
func Load(root string, profile string) (Config, error) {
	if profile == "" {
		profile = envOr("APP_ENV", "development")
	}

	values := map[string]string{}
	for _, name := range []string{".env." + profile, ".env.local"} {
		if err := readInto(values, filepath.Join(root, name)); err != nil {
			return Config{}, err
		}
	}
	for _, key := range []string{"ENVIRONMENT", "LOG_LEVEL", "VERBOSE_ERRORS", "CACHE_SECONDS"} {
		if v, ok := os.LookupEnv(key); ok {
			values[key] = v
		}
	}

	cfg := Config{
		Environment:   values["ENVIRONMENT"],
		LogLevel:      values["LOG_LEVEL"],
		VerboseErrors: values["VERBOSE_ERRORS"] == "true",
	}
	if !contains(Environments, cfg.Environment) {
		return Config{}, fmt.Errorf("config: unknown ENVIRONMENT %q (want one of %v)", cfg.Environment, Environments)
	}
	if cfg.LogLevel == "" {
		return Config{}, fmt.Errorf("config: LOG_LEVEL is missing from profile %q", profile)
	}
	seconds, err := strconv.Atoi(values["CACHE_SECONDS"])
	if err != nil || seconds < 0 {
		return Config{}, fmt.Errorf("config: CACHE_SECONDS must be a non-negative integer, got %q", values["CACHE_SECONDS"])
	}
	cfg.CacheSeconds = seconds

	return cfg, nil
}

// readInto parses KEY=VALUE lines. A missing file is not an error: only the
// selected profile has to exist, and .env.local usually does not.
func readInto(values map[string]string, path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		values[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(value), `"'`)
	}
	return scanner.Err()
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func contains(list []string, s string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}
