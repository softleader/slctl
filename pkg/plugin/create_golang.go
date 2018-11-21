package plugin

import (
	"path/filepath"
)

const (
	MainFileName     = "main.go"
	MakefileFileName = "Makefile"
)

const defaultMain = `package main
import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func main() {
	cmd := &cobra.Command{
		Use:   "{{.Name}}",
		Short: "{{.Usage}}",
		Long:  "{{.Description}}",
		RunE: func(cmd *cobra.Command, args []string) error {
			// use os.LookupEnv to retrieve the specific value of the environment (e.g. os.LookupEnv("SL_TOKEN"))
			for _, env := range os.Environ() {
				if strings.HasPrefix(env, "SL_") {
					fmt.Println(env)
				}
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
VERSION := $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' plugin.yaml)
DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
LDFLAGS := "-X main.version=${VERSION}"
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

.PHONY: dist
dist:
	go get -u github.com/inconshreveable/mousetrap
	mkdir -p $(BUILD)
	mkdir -p $(DIST)
	sed -E 's/(version: )"(.+)"/\1"$(VERSION)"/g' plugin.yaml > $(BUILD)/plugin.yaml
	GOOS=linux GOARCH=amd64 go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS) -a -tags netgo
	tar -C $(BUILD) -zcvf $(DIST)/$(BINARY)-linux-$(VERSION).tgz $(BINARY) plugin.yaml
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS) -a -tags netgo
	tar -C $(BUILD) -zcvf $(DIST)/$(BINARY)-macos-$(VERSION).tgz $(BINARY) plugin.yaml
	GOOS=windows GOARCH=amd64 go build -o $(BUILD)/$(BINARY).exe -ldflags $(LDFLAGS) -a -tags netgo
	tar -C $(BUILD) -llzcvf $(DIST)/$(BINARY)-windows-$(VERSION).tgz $(BINARY).exe plugin.yaml

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

type golang struct{}

func (c golang) files(plugin *Metadata, pluginDir string) []file {
	return []file{
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
}
