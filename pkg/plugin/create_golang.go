package plugin

import (
	"path/filepath"
)

const golangMain = `package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
	"io"
)

type {{.Name|lowerCamel}}Cmd struct {
	out     io.Writer	
	offline bool
	verbose bool
	token   string
}

func main() {
	c := {{.Name|lowerCamel}}Cmd{}
	cmd := &cobra.Command{
		Use:   "{{.Name}}",
		Short: "{{.Usage}}",
		Long:  "{{.Description}}",
		RunE: func(cmd *cobra.Command, args []string) error {
			c.token = os.ExpandEnv(c.token)
			return c.run()
		},
	}
	
	c.out = cmd.OutOrStdout()
	c.offline, _ = strconv.ParseBool(os.Getenv("SL_OFFLINE"))
	c.verbose, _ = strconv.ParseBool(os.Getenv("SL_VERBOSE"))

	f := cmd.Flags()
	f.BoolVarP(&c.offline, "offline", "o", c.offline, "work offline, Overrides $SL_OFFLINE")
	f.BoolVarP(&c.verbose, "verbose", "v", c.verbose, "enable verbose output, Overrides $SL_VERBOSE")
	f.StringVar(&c.token, "token", "$SL_TOKEN", "github access token. Overrides $SL_TOKEN")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func (c *{{.Name|lowerCamel}}Cmd) run() error {
	// use os.LookupEnv to retrieve the specific value of the environment (e.g. os.LookupEnv("SL_TOKEN"))
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "SL_") {
			fmt.Println(env)
		}
	}
	fmt.Printf("%+v\n", c)
	return nil
}
`

const golangVersion = `package main

import (
	"strings"
)

const (
	unreleased  = "unreleased"
)

var version string

func ver() string {
	if v := strings.TrimSpace(version); v != "" {
		return v
	} else {
		return unreleased
	}
}
`

const golangMakefile = `SL_HOME ?= $(shell slctl home)
SL_PLUGIN_DIR ?= $(SL_HOME)/plugins/{{.Name}}/
METADATA := metadata.yaml
HAS_GLIDE := $(shell command -v glide;)
VERSION := $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' $(METADATA))
DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
LDFLAGS := "-X main.version=${VERSION}"
BINARY := {{.Name}}

.PHONY: install
install: bootstrap test build
	mkdir -p $(SL_PLUGIN_DIR)
	cp $(BUILD)/$(BINARY) $(SL_PLUGIN_DIR)
	cp $(METADATA) $(SL_PLUGIN_DIR)

.PHONY: test
test:
	go test ./... -v

.PHONY: build
build: clean bootstrap
	mkdir -p $(BUILD)
	cp $(METADATA) $(BUILD)
	go build -o $(BUILD)/$(BINARY)

.PHONY: dist
dist:
	go get -u github.com/inconshreveable/mousetrap
	mkdir -p $(BUILD)
	mkdir -p $(DIST)
	sed -E 's/^(version: )(.+)/\1$(VERSION)/g' $(METADATA) > $(BUILD)/$(METADATA)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS) -a -tags netgo
	tar -C $(BUILD) -zcvf $(DIST)/$(BINARY)-linux-$(VERSION).tgz $(BINARY) $(METADATA)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS) -a -tags netgo
	tar -C $(BUILD) -zcvf $(DIST)/$(BINARY)-darwin-$(VERSION).tgz $(BINARY) $(METADATA)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD)/$(BINARY).exe -ldflags $(LDFLAGS) -a -tags netgo
	tar -C $(BUILD) -llzcvf $(DIST)/$(BINARY)-windows-$(VERSION).tgz $(BINARY).exe $(METADATA)

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

func (c golang) exec(plugin *Metadata) Commands {
	command := "$SL_PLUGIN_DIR/" + plugin.Name
	return Commands{
		Command: command,
		Platform: []Platform{
			{Os: "darwin", Command: command,},
			{Os: "windows", Command: command,},
		},
	}
}

func (c golang) hook(plugin *Metadata) Commands {
	return Commands{
		Command: "echo hello " + plugin.Name,
	}
}

func (c golang) files(plugin *Metadata, pdir string) []file {
	return []file{
		tpl{
			path:     filepath.Join(pdir, "main.go"),
			in:       plugin,
			template: golangMain,
		},
		tpl{
			path:     filepath.Join(pdir, "version.go"),
			in:       plugin,
			template: golangVersion,
		},
		tpl{
			path:     filepath.Join(pdir, "Makefile"),
			in:       plugin,
			template: golangMakefile,
		},
	}
}
