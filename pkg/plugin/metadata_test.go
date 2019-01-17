package plugin

import "testing"

func TestMetadata_IsVersionGreaterThan(t *testing.T) {
	a := Metadata{
		Version: "1.0.0",
	}
	b := Metadata{
		Version: "1.1.0",
	}
	isGreater, err := b.IsVersionGreaterThan(&a)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if !isGreater {
		t.Error("b should greater than a")
	}

	b = Metadata{
		Version: "0.0.0",
	}
	isGreater, err = b.IsVersionGreaterThan(&a)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if isGreater {
		t.Error("b should not greater than a")
	}
}

func TestMetadata_IsVersionLegal(t *testing.T) {
	v := "0.0.0"
	if !(&Metadata{Version: v}).IsVersionLegal() {
		t.Errorf("%s should legal", v)
	}
	v = "0.0.hello"
	if (&Metadata{Version: v}).IsVersionLegal() {
		t.Errorf("%s should hot legal", v)
	}
}
