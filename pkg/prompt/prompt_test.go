package prompt

import (
	"bytes"
	"strings"
	"testing"
)

func TestYesNoQuestionFrom(t *testing.T) {
	out := &bytes.Buffer{}
	
	// Yes
	inY := strings.NewReader("y\n")
	if !YesNoQuestionFrom(inY, out, "Continue?") {
		t.Error("expected true for y")
	}

	// No
	inN := strings.NewReader("n\n")
	if YesNoQuestionFrom(inN, out, "Continue?") {
		t.Error("expected false for n")
	}

	// Retry then Yes
	inRetry := strings.NewReader("invalid\nyes\n")
	if !YesNoQuestionFrom(inRetry, out, "Continue?") {
		t.Error("expected true for yes after retry")
	}
}
