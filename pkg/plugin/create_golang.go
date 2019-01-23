package plugin

import (
	"fmt"
	"path/filepath"
	"strings"
)

const golangMain = `package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var (
	// 在包版時會動態指定 version 及 commit
	version, commit string
	metadata        *release.Metadata

	// global flags
	verbose, offline bool
	token            string
)

func main() {
	cobra.OnInitialize(
		initMetadata,
		initGlobalFlags,
		initFlags,
	)

	cmd := &cobra.Command{
		Use:   "{{.Name}}",
		Short: "{{.Usage}}",
		Long:  "{{.Description}}",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// remove the check if the plugin can run in offline mode
			if offline {
				return fmt.Errorf("can not run the command in offline mode")
			}
			logrus.SetOutput(cmd.OutOrStdout())
			logrus.SetFormatter(&formatter.PlainFormatter{})

			// use os.LookupEnv to retrieve the specific value of the environment (e.g. os.LookupEnv("SL_TOKEN"))
			for _, env := range os.Environ() {
				if strings.HasPrefix(env, "SL_") {
					logrus.Println(env)
				}
			}
			return nil
		},
	}

	f := cmd.PersistentFlags()
	f.BoolVar(&offline, "offline", offline, "work offline, Overrides $SL_OFFLINE")
	f.BoolVarP(&verbose, "verbose", "v", verbose, "enable verbose output, Overrides $SL_VERBOSE")
	f.StringVar(&token, "token", token, "github access token. Overrides $SL_TOKEN")

	cmd.AddCommand(
		newVersionCmd(),
	)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// initMetadata 準備 app 的 release 資訊
func initMetadata() {
	metadata = release.NewMetadata(version, commit)
}

// initGlobalFlags 準備 app 的 global flags 預設值
func initGlobalFlags() {
	offline, _ = strconv.ParseBool(os.Getenv("SL_OFFLINE"))
	verbose, _ = strconv.ParseBool(os.Getenv("SL_VERBOSE"))
	token = os.Getenv("SL_TOKEN")
}
`

const golangVersion = `package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	var full bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print {{.Name}} version",
		Long:  "print {{.Name}} version",
		RunE: func(cmd *cobra.Command, args []string) error {
			if full {
				logrus.Infoln(metadata.FullString())
			} else {
				logrus.Infoln(metadata.String())
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.BoolVar(&full, "full", false, "print full version number and commit hash")

	return cmd
}
`

const golangPkgFormatter = `package formatter

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

var ln = fmt.Sprintln()

// PlainFormatter 代表什麼都不 format 的 formatter
type PlainFormatter struct {
}

// Format 將傳入的 entry 轉換成要寫 log 的文字
func (f *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buf *bytes.Buffer
	if entry.Buffer != nil {
		buf = entry.Buffer
	} else {
		buf = &bytes.Buffer{}
	}
	if entry.Message != "" {
		buf.WriteString(entry.Message)
	}
	if !strings.HasSuffix(entry.Message, ln) {
		buf.WriteString(ln)
	}
	return buf.Bytes(), nil
}
`

const golangPkgFormatterTest = `package formatter

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestPlainFormatter_Format(t *testing.T) {
	log := logrus.New()
	log.SetFormatter(&PlainFormatter{})
	b := bytes.NewBuffer(nil)
	log.SetOutput(b)
	log.Println("123456789")
	if b.String() != "123456789\n" {
		t.Error("out should be 123456789\n")
	}
}
`
const golangPkgRelease = `package release

import (
	"fmt"
	"strings"
)

const (
	unreleased = "unreleased"
	unknown    = "unknown"
)

// Metadata 代表此 app 的 release 資訊
type Metadata struct {
	GitVersion string
	GitCommit  string
}

// NewMetadata 產生一個 app 的 release 資訊
func NewMetadata(version, commit string) (b *Metadata) {
	b = &Metadata{
		GitVersion: unreleased,
		GitCommit:  unknown,
	}
	if version = strings.TrimSpace(version); version != "" {
		b.GitVersion = version
	}
	if commit = strings.TrimSpace(commit); commit != "" {
		b.GitCommit = commit
	}
	return
}

func (b *Metadata) String() string {
	trunc := 7
	if len := len(b.GitCommit); len < 7 {
		trunc = len
	}
	return fmt.Sprintf("%s+%s", b.GitVersion, b.GitCommit[:trunc])
}

// FullString 回傳完整的 release 資訊
func (b *Metadata) FullString() string {
	return fmt.Sprintf("%#v", b)
}
`

const golangPkgReleaseTest = `package release

import (
	"fmt"
	"testing"
)

func TestMetadata_String(t *testing.T) {
	commit := "none"
	expected := fmt.Sprintf("%s+%s", unreleased, commit)
	if v := NewMetadata(unreleased, commit).String(); v != expected {
		t.Errorf("expected to see %q, but got %q", expected, v)
	}
	commit = "asdfbngfdseqw2314rtygfsda"
	expected = fmt.Sprintf("%s+%s", unreleased, commit[:7])
	if v := NewMetadata(unreleased, commit).String(); v != expected {
		t.Errorf("expected to see %q, but got %q", expected, v)
	}
}
`

const golangMakefile = `HAS_GOLINT := $(shell command -v golint;)
SL_HOME ?= $(shell slctl home)
SL_PLUGIN_DIR ?= $(SL_HOME)/plugins/{{.Name}}/
METADATA := metadata.yaml
VERSION := $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' $(METADATA))
DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
LDFLAGS := "-X main.version=${VERSION}"
BINARY := {{.Name}}
MAIN := ./cmd/{{.Name|lower}}

.PHONY: install
install: bootstrap test build
	mkdir -p $(SL_PLUGIN_DIR)
	cp $(BUILD)/$(BINARY) $(SL_PLUGIN_DIR)
	cp $(METADATA) $(SL_PLUGIN_DIR)

.PHONY: test
test: golint
	go test ./... -v

.PHONY: gofmt
gofmt:
	gofmt -s -w .

.PHONY: golint
golint: gofmt
ifndef HAS_GOLINT
	go get -u golang.org/x/lint/golint
endif
	golint -set_exit_status ./cmd/...
	golint -set_exit_status ./pkg/...

.PHONY: build
build: clean bootstrap
	mkdir -p $(BUILD)
	cp $(METADATA) $(BUILD)
	go build -o $(BUILD)/$(BINARY) $(MAIN)

.PHONY: dist
dist:
	go get -u github.com/inconshreveable/mousetrap
	mkdir -p $(BUILD)
	mkdir -p $(DIST)
	sed -E 's/^(version: )(.+)/\1$(VERSION)/g' $(METADATA) > $(BUILD)/$(METADATA)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS) -a -tags netgo $(MAIN)
	tar -C $(BUILD) -zcvf $(DIST)/$(BINARY)-linux-$(VERSION).tgz $(BINARY) $(METADATA)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS) -a -tags netgo $(MAIN)
	tar -C $(BUILD) -zcvf $(DIST)/$(BINARY)-darwin-$(VERSION).tgz $(BINARY) $(METADATA)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD)/$(BINARY).exe -ldflags $(LDFLAGS) -a -tags netgo $(MAIN)
	tar -C $(BUILD) -llzcvf $(DIST)/$(BINARY)-windows-$(VERSION).tgz $(BINARY).exe $(METADATA)

.PHONY: bootstrap
bootstrap:
ifeq (,$(wildcard ./go.mod))
	go mod init {{.Name}}
endif
	go mod download

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
			{Os: "darwin", Command: command},
			{Os: "windows", Command: command},
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
			path:     filepath.Join(pdir, "cmd", cmd, "version.go"),
			in:       plugin,
			template: golangVersion,
		},
		tpl{
			path:     filepath.Join(pdir, "pkg", "formatter", "formatter.go"),
			in:       plugin,
			template: golangPkgFormatter,
		},
		tpl{
			path:     filepath.Join(pdir, "pkg", "formatter", "formatter_test.go"),
			in:       plugin,
			template: golangPkgFormatterTest,
		},
		tpl{
			path:     filepath.Join(pdir, "pkg", "release", "release.go"),
			in:       plugin,
			template: golangPkgRelease,
		},
		tpl{
			path:     filepath.Join(pdir, "pkg", "release", "release_test.go"),
			in:       plugin,
			template: golangPkgReleaseTest,
		},
		tpl{
			path:     filepath.Join(pdir, "Makefile"),
			in:       plugin,
			template: golangMakefile,
		},
	}
}
