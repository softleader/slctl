package cmd

import (
	"bytes"
	"fmt"
	"testing"
)

func TestConfirmToken(t *testing.T) {
	settings.Debug = true
	b := bytes.NewBuffer(nil)
	token := "997f19253fccc351bfcf4cf1622f494f7708522a"
	var err error
	var name string

	if name, err = confirmToken(token, b); err != nil {
		t.Error(err)
	}
	if name == "" {
		t.Errorf("name should not be empty")
	}

	fmt.Printf("Hello, %s!\n", name)
}
