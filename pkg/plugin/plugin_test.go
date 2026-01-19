package plugin

import (
	"os"
	"testing"
)

func TestFromGitHub(t *testing.T) {
	tests := []struct {
		source string
		want   bool
	}{
		{"github.com/softleader/slctl", true},
		{"/home/user/plugin", false},
	}

	for _, tt := range tests {
		p := &Plugin{Source: tt.source}
		if got := p.FromGitHub(); got != tt.want {
			t.Errorf("FromGitHub(%q) = %v, want %v", tt.source, got, tt.want)
		}
	}
}

func TestPrepareCommand(t *testing.T) {
	p := &Plugin{
		Metadata: &Metadata{
			IgnoreGlobalFlags: false,
		},
	}

	os.Setenv("VAR", "value")
	defer os.Unsetenv("VAR")

	main, argv, err := p.PrepareCommand("echo $VAR", []string{"--flag"})
	if err != nil {
		t.Fatal(err)
	}
	if main != "echo" {
		t.Errorf("expected echo, got %s", main)
	}
	if argv[0] != "value" {
		t.Errorf("expected value, got %s", argv[0])
	}
	if argv[1] != "--flag" {
		t.Errorf("expected --flag, got %s", argv[1])
	}

	// Test IgnoreGlobalFlags
	p.Metadata.IgnoreGlobalFlags = true
	main, argv, err = p.PrepareCommand("echo $VAR", []string{"--flag"})
	if len(argv) != 1 {
		t.Errorf("expected 1 arg, got %d", len(argv))
	}
}
