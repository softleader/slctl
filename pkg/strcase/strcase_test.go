package strcase

import (
	"testing"
)

func TestToLowerCamel(t *testing.T) {
	str := "DSfglctl-foo_sadfghewq_qwewrqwdsfaSdas"
	expected := "dSfglctlFooSadfghewqQwewrqwdsfaSdas"
	if actual := ToLowerCamel(str); actual != expected {
		t.Error("expected to be " + expected)
	}
}
