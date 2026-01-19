package ver

import (
	"testing"
)

func TestRevision_String(t *testing.T) {
	r := Revision("1.2.3")
	if r.String() != "1.2.3" {
		t.Errorf("expected 1.2.3, got %s", r.String())
	}
}

func TestRevision_IsGreaterThan(t *testing.T) {
	r := Revision("1.2.3")

	tests := []struct {
		other   string
		want    bool
		wantErr bool
	}{
		{"1.2.2", true, false},
		{"1.2.3", false, false},
		{"1.2.4", false, false},
		{"invalid", false, true},
	}

	for _, tt := range tests {
		got, err := r.IsGreaterThan(tt.other)
		if (err != nil) != tt.wantErr {
			t.Errorf("IsGreaterThan(%q) error = %v, wantErr %v", tt.other, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("IsGreaterThan(%q) = %v, want %v", tt.other, got, tt.want)
		}
	}

	// Invalid current revision
	r2 := Revision("invalid")
	_, err := r2.IsGreaterThan("1.0.0")
	if err == nil {
		t.Error("expected error for invalid current revision")
	}
}
