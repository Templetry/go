package config

import (
	"path/filepath"
	"testing"
)

// The profiles live at the module root, three levels above this package.
func root(t *testing.T) string {
	t.Helper()
	return filepath.Join("..", "..")
}

// Proves the profiles are wired, not decorative: each one is loaded and read
// back, so renaming a file or a key breaks here rather than in whatever
// environment happened to depend on it.
func TestEachProfileDeclaresItsOwnName(t *testing.T) {
	for _, name := range Environments {
		cfg, err := Load(root(t), name)
		if err != nil {
			t.Fatalf("loading %s: %v", name, err)
		}
		if cfg.Environment != name {
			t.Errorf("%s: Environment = %q, want %q", name, cfg.Environment, name)
		}
	}
}

func TestProfilesDifferWhereItMatters(t *testing.T) {
	cases := []struct {
		profile string
		verbose bool
		cache   int
	}{
		{"development", true, 0},
		// Staging exists to be production-like while still debuggable, so it
		// is the one profile whose values must not equal either neighbour's.
		{"staging", true, 30},
		{"production", false, 300},
	}
	for _, c := range cases {
		cfg, err := Load(root(t), c.profile)
		if err != nil {
			t.Fatalf("loading %s: %v", c.profile, err)
		}
		if cfg.VerboseErrors != c.verbose || cfg.CacheSeconds != c.cache {
			t.Errorf("%s: got verbose=%v cache=%d, want %v and %d",
				c.profile, cfg.VerboseErrors, cfg.CacheSeconds, c.verbose, c.cache)
		}
	}
}

func TestAppEnvSelectsTheProfile(t *testing.T) {
	t.Setenv("APP_ENV", "production")

	cfg, err := Load(root(t), "")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Environment != "production" {
		t.Errorf("Environment = %q, want production", cfg.Environment)
	}
}

// This is the property that makes containers work: the platform sets real
// variables and no profile file has to ship with the image.
func TestARealEnvironmentVariableWinsOverTheFile(t *testing.T) {
	t.Setenv("CACHE_SECONDS", "7")

	cfg, err := Load(root(t), "production")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.CacheSeconds != 7 {
		t.Errorf("CacheSeconds = %d, want 7", cfg.CacheSeconds)
	}
}

func TestAnUnknownProfileFailsLoudly(t *testing.T) {
	if _, err := Load(root(t), "qa"); err == nil {
		t.Error("loading a profile that does not exist should be an error")
	}
}
