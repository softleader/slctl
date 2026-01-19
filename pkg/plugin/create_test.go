package plugin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreate(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "sl-plugin-create")
	defer os.RemoveAll(tempDir)

	plugin := &Metadata{
		Name:    "test-go-plugin",
		Version: "0.1.0",
	}

	path, err := Create("golang", plugin, filepath.Join(tempDir, "go-plugin"))
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(path, MetadataFileName)); err != nil {
		t.Errorf("expected metadata.yaml to exist, got %v", err)
	}

	// Java
	path, err = Create("java", plugin, filepath.Join(tempDir, "java-plugin"))
	if err != nil {
		t.Fatalf("Create java failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(path, MetadataFileName)); err != nil {
		t.Errorf("expected metadata.yaml to exist, got %v", err)
	}

	// NodeJS
	path, err = Create("nodejs", plugin, filepath.Join(tempDir, "node-plugin"))
	if err != nil {
		t.Fatalf("Create nodejs failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(path, MetadataFileName)); err != nil {
		t.Errorf("expected metadata.yaml to exist, got %v", err)
	}
}

func TestMarshal(t *testing.T) {
	m := marshal{
		path: "test.yaml",
		in:   map[string]string{"foo": "bar"},
	}
	if m.filepath() != "test.yaml" {
		t.Errorf("expected test.yaml, got %s", m.filepath())
	}
	content, err := m.content()
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "foo: bar\n" {
		t.Errorf("expected foo: bar\\n, got %s", string(content))
	}
}

func TestTpl(t *testing.T) {
	type data struct {
		Name string
	}
	u := tpl{
		path:     "test.txt",
		in:       data{Name: "World"},
		template: "Hello {{.Name | upper}}",
	}
	content, err := u.content()
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "Hello WORLD" {
		t.Errorf("expected Hello WORLD, got %s", string(content))
	}
}
