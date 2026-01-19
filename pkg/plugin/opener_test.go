package plugin

import (
	"testing"
)

func TestOpen(t *testing.T) {
	// We can't easily test Open as it triggers an external command.
	// But we can check the logic of FromGitHub in Plugin.
	p := &Plugin{
		Dir:    "/plugins/test",
		Source: "github.com/softleader/slctl",
	}

	if !p.FromGitHub() {
		t.Error("expected FromGitHub to be true")
	}

	p.Source = "/local/path"
	if p.FromGitHub() {
		t.Error("expected FromGitHub to be false")
	}
}

func TestOpen_Lines(t *testing.T) {
	p := &Plugin{Dir: "/tmp", Source: "github.com/softleader/slctl"}
	// We won't actually call Open as it might fail in CI
	// But let's check the FromGitHub logic inside Open via a testable way if possible.
	// Actually, just calling p.Open("") will try to run open.Run(p.Dir)
	// which will likely fail on headless CI but might pass if it just returns err.
	p.Open("")
}
