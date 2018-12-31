package plugin

import (
	"fmt"
	"path/filepath"
	"strings"
)

const golangMain = `package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	unreleased  = "unreleased"
	unknown     = "unknown"
)

var (
	version string
	commit  string
	date    string
)

type Version struct {
	GitVersion string
	GitCommit  string
	BuildDate  string
}

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

	cmd.AddCommand(&cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(c.out, "%+v\n", ver())
		},
	})

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func (c *fooCmd) run() error {
	// use os.LookupEnv to retrieve the specific value of the environment (e.g. os.LookupEnv("SL_TOKEN"))
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "SL_") {
			fmt.Fprintln(c.out, env)
		}
	}
	fmt.Fprintf(c.out, "%+v\n", c)
	return nil
}

func ver() *Version {
	if version = strings.TrimSpace(version); version == "" {
		version = unreleased
	}
	if commit = strings.TrimSpace(commit); commit == "" {
		commit = unknown
	}
	if date = strings.TrimSpace(date); date == "" {
		date = unknown
	}
	return &Version{
		GitVersion: version,
		GitCommit:  commit,
		BuildDate:  date,
	}
}
`

const golangGoReleaser = `before:
  hooks:
    - go mod download
builds:
  - main: ./cmd/{{.Name|lower}}
    env:
      - CGO_ENABLED=0
    goos:
      - darwin
      - linux
      - windows
    goarch:
      - amd64
dist: _dist
archive:
  replacements:
    darwin: darwin
    linux: linux
    window: windows
  files:
    - metadata.yaml
checksum:
  name_template: 'checksums.txt'
snapshot:
  name_template: "{{.Name|lower}}-SNAPSHOT"
changelog:
  sort: asc
  filters:
    exclude:
      - '^docs:'
      - '^test:'`

const golangMakefile = `BUILD := $(CURDIR)/_build
BINARY := {{.Name}}
MAIN := ./cmd/{{.Name|lower}}
BUILD := $(CURDIR)/_build

.PHONY: test
test:
	go test ./... -v

.PHONY: build
build:
	go build -o $(BUILD)/$(BINARY) $(MAIN)

.PHONY: dist
dist: bootstrap
	goreleaser release --snapshot --rm-dist

.PHONY: bootstrap
bootstrap:
ifeq (,$(wildcard ./go.mod))
	go mod init {{.Name}}
endif
	go mod download

.PHONY: clean
clean:
	rm -rf _*
	rm -f /usr/local/bin/$(BINARY)
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
	cmd := strings.ToLower(plugin.Name)
	return []file{
		tpl{
			path:     filepath.Join(pdir, "cmd", cmd, fmt.Sprintf("%s.go", cmd)),
			in:       plugin,
			template: golangMain,
		},
		tpl{
			path:     filepath.Join(pdir, ".goreleaser.yml"),
			in:       plugin,
			template: golangGoReleaser,
		},
		tpl{
			path:     filepath.Join(pdir, "Makefile"),
			in:       plugin,
			template: golangMakefile,
		},
	}
}
