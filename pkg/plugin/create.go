package plugin

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	PluginFileName   = "plugin.yaml"
	MainFileName     = "main.go"
	MakefileFileName = "Makefile"
)

const defaultMain = `package main
import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	cmd := &cobra.Command{
		Use:   "{{.Name}}",
		Short: "{{.Usage}}",
		Long:  "{{.Description}}",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				fmt.Println(arg)
			}
			return nil
		},
	}
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
`

const defaultMakefile = `SL_HOME ?= $(shell slctl home)
SL_PLUGIN_DIR ?= $(SL_HOME)/plugins/{{.Name}}/
HAS_GLIDE := $(shell command -v glide;)
VERSION := ""
BUILD := $(CURDIR)/_build
BINARY := {{.Name}}

.PHONY: install
install: bootstrap test build
	mkdir -p $(SL_PLUGIN_DIR)
	cp $(BUILD)/$(BINARY) $(SL_PLUGIN_DIR)
	cp plugin.yaml $(SL_PLUGIN_DIR)

.PHONY: test
test:
	go test -v

.PHONY: build
build: clean bootstrap
	mkdir -p $(BUILD)
	cp plugin.yaml $(BUILD)
	go build -o $(BUILD)/$(BINARY)

.PHONY: bootstrap
bootstrap:
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif
ifeq (,$(wildcard ./glide.yaml))
	glide init --non-interactive
endif
	glide install --strip-vendor

.PHONY: clean
clean:
	rm -rf _*
`

func Create(plugin *Metadata, dir string) (string, error) {
	path, err := filepath.Abs(dir)
	if err != nil {
		return path, err
	}

	if fi, err := os.Stat(path); err != nil {
		return path, err
	} else if !fi.IsDir() {
		return path, fmt.Errorf("no such directory %s", path)
	}

	pluginDir := filepath.Join(path, plugin.Name)
	if fi, err := os.Stat(pluginDir); err == nil && !fi.IsDir() {
		return pluginDir, fmt.Errorf("file %s already exists and is not a directory", pluginDir)
	}
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return pluginDir, err
	}

	units := []unit{
		marshal{
			path: filepath.Join(pluginDir, PluginFileName),
			in:   plugin,
		},
		compile{
			path:     filepath.Join(pluginDir, MainFileName),
			in:       plugin,
			template: defaultMain,
		},
		compile{
			path:     filepath.Join(pluginDir, MakefileFileName),
			in:       plugin,
			template: defaultMakefile,
		},
	}

	for _, unit := range units {
		if err := save(unit); err != nil {
			return pluginDir, err
		}
	}
	return pluginDir, nil
}

type unit interface {
	filename() string
	content() ([]byte, error)
}

type marshal struct {
	path string
	in   interface{}
}

func (u marshal) filename() string {
	return u.path
}

func (u marshal) content() ([]byte, error) {
	return yaml.Marshal(u.in)
}

type compile struct {
	path     string
	in       interface{}
	template string
}

func (u compile) filename() string {
	return u.path
}

func (u compile) content() ([]byte, error) {
	var buf bytes.Buffer
	parsed := template.Must(template.New("").Parse(u.template))
	if err := parsed.Execute(&buf, u.in); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func save(unit unit) (err error) {
	out, err := unit.content()
	if err != nil {
		return
	}
	return ioutil.WriteFile(unit.filename(), out, 0644)
}
