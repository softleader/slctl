package plugin

import (
	"runtime"
	"testing"
)

func TestGetCommand(t *testing.T) {
	// Case 1: Generic command only
	c1 := &Commands{Command: "universal"}
	cmd, err := c1.GetCommand()
	if err != nil {
		t.Error(err)
	}
	if cmd != "universal" {
		t.Errorf("expected universal, got %s", cmd)
	}

	// Case 2: Matching OS and Arch
	c2 := &Commands{
		Command: "universal",
		Platform: []Platform{
			{Os: runtime.GOOS, Arch: runtime.GOARCH, Command: "exact-match"},
		},
	}
	cmd, err = c2.GetCommand()
	if err != nil {
		t.Error(err)
	}
	if cmd != "exact-match" {
		t.Errorf("expected exact-match, got %s", cmd)
	}

	// Case 3: Matching OS only
	c3 := &Commands{
		Command: "universal",
		Platform: []Platform{
			{Os: runtime.GOOS, Arch: "not-my-arch", Command: "os-match"},
		},
	}
	cmd, err = c3.GetCommand()
	if err != nil {
		t.Error(err)
	}
	if cmd != "os-match" {
		t.Errorf("expected os-match, got %s", cmd)
	}

	// Case 4: No command found
	c4 := &Commands{}
	_, err = c4.GetCommand()
	if err == nil {
		t.Error("expected error for no command found")
	}
}

func TestErrNoCommandFound_Error(t *testing.T) {
	err := &ErrNoCommandFound{s: "test error"}
	if err.Error() != "test error" {
		t.Errorf("expected test error, got %s", err.Error())
	}
}
