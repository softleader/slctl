package main

import (
	"fmt"
	"testing"
)

func TestVersion_String(t *testing.T) {
	commit = "none"
	expected := fmt.Sprintf("%s+%s", unreleased, commit)
	if v := ver().String(); v != expected {
		t.Errorf("expected to see %q, but got %q", expected, v)
	}
	commit = "asdfbngfdseqw2314rtygfsda"
	expected = fmt.Sprintf("%s+%s", unreleased, commit[:7])
	if v := ver().String(); v != expected {
		t.Errorf("expected to see %q, but got %q", expected, v)
	}
}
