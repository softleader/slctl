package formatter

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestPlainFormatter_Format(t *testing.T) {
	f := &PlainFormatter{}
	entry := &logrus.Entry{Message: "hello"}
	
got, err := f.Format(entry)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello\n" {
		t.Errorf("expected hello\\n, got %q", string(got))
	}

	entryWithBuf := &logrus.Entry{Message: "world"}
	got, _ = f.Format(entryWithBuf)
	if string(got) != "world\n" {
		t.Errorf("expected world\\n, got %q", string(got))
	}
}

