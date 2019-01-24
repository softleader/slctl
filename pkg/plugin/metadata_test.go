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

	a = Metadata{
		Version: "2.5.1",
	}
	b = Metadata{
		Version: "2.6.4",
	}
	isGreater, err = a.IsVersionGreaterThan(&b)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if isGreater {
		t.Error("a should not greater than b")
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
