package greet

import "testing"

func TestGreeting(t *testing.T) {
	cases := []struct{ in, want string }{
		{"world", "Hello, world!"},
		{"", "Hello, world!"},
		{"Go", "Hello, Go!"},
	}
	for _, c := range cases {
		if got := Greeting(c.in); got != c.want {
			t.Errorf("Greeting(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
