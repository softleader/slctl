package plugin

import (
	"path/filepath"
)

const nodejsIndex = `const {Command, flags} = require('@oclif/command')

class {{.Name|title}}Command extends Command {
  async run() {
    const {flags} = this.parse({{.Name|title}}Command)
    const name = flags.name || 'world'
    Object.keys(process.env)
      .filter(key => key.startsWith("SL_"))
      .forEach(key => this.log(key + "=" + process.env[key]));
  }
}

{{.Name|title}}Command.description = 'Describe the command here'

{{.Name|title}}Command.flags = {
  // add --version flag to show CLI version
  version: flags.version({char: 'v'}),
  // add --help flag to show CLI version
  help: flags.help({char: 'h'}),
  name: flags.string({char: 'n', description: 'name to print'}),
}

module.exports = {{.Name|title}}Command
`

const nodejsPackageJson = `{
  "name": "{{.Name}}",
  "version": "{{.Version}}",
  "author": "@softleader",
  "bin": {
    "{{.Name}}": "./bin/run"
  },
  "dependencies": {
    "@oclif/command": "^1.5.6",
    "@oclif/config": "^1.9.0",
    "@oclif/plugin-help": "^2.1.4"
  },
  "engines": {
    "node": ">=8.0.0"
  },
  "files": [
    "/bin",
    "/src"
  ],
  "keywords": [
    "oclif"
  ],
  "license": "MIT",
  "main": "src/index.js",
  "oclif": {
    "bin": "{{.Name}}"
  },
  "scripts": {
    "test": "echo NO TESTS"
  }
}
`

const nodejsRun = `#!/usr/bin/env node

require('..').run()
.catch(require('@oclif/errors/handle'))`

const nodejsRunCmd = `@echo off

node "%~dp0\run" %*`

const nodejsMakefile = `SL_HOME ?= $(shell slctl home)
SL_PLUGIN_DIR ?= $(SL_HOME)/plugins/{{.Name}}/
METADATA := metadata.yaml
HAS_NODE := $(shell command -v node;)
VERSION := $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' $(METADATA))
DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
MODULES := $(CURDIR)/node_modules
BIN := $(CURDIR)/bin
SRC := $(CURDIR)/src

.PHONY: install
install: build
	mkdir -p $(SL_PLUGIN_DIR)
	cp -r $(BUILD)/* $(SL_PLUGIN_DIR)

.PHONY: build
build: clean bootstrap
	mkdir -p $(BUILD)
	sed -E 's/(version: )"(.+)"/\1"$(VERSION)"/g' $(METADATA) > $(BUILD)/$(METADATA)
	cp package.json $(BUILD)
	cp -r node_modules $(BUILD)
	cp -r src $(BUILD)
	cp -r bin $(BUILD)

.PHONY: dist
dist: build
	mkdir -p $(DIST)
	tar -C $(BUILD) -zcvf $(DIST)/{{.Name}}-$(VERSION).tgz $(METADATA) package.json node_modules src bin

.PHONY: bootstrap
bootstrap:
ifndef HAS_NODE
	$(error You must install Node.js)
endif
	npm install

.PHONY: clean
clean:
	rm -rf _*
`

type nodejs struct{}

func (c nodejs) exec(plugin *Metadata) Commands {
	command := "node $SL_PLUGIN_DIR/bin/run"
	return Commands{
		Command: command,
		Platform: []Platform{
			{Os: "darwin", Command: command,},
			{Os: "windows", Command: command,},
		},
	}
}

func (c nodejs) hook(plugin *Metadata) Commands {
	return Commands{
		Command: "echo hello " + plugin.Name,
	}
}

func (c nodejs) files(plugin *Metadata, pdir string) []file {
	return []file{
		tpl{
			path:     filepath.Join(pdir, "src", "index.js"),
			in:       plugin,
			template: nodejsIndex,
		},
		tpl{
			path:     filepath.Join(pdir, "package.json"),
			in:       plugin,
			template: nodejsPackageJson,
		},
		tpl{
			path:     filepath.Join(pdir, "bin", "run"),
			in:       plugin,
			template: nodejsRun,
		},
		tpl{
			path:     filepath.Join(pdir, "bin", "run.cmd"),
			in:       plugin,
			template: nodejsRunCmd,
		},
		tpl{
			path:     filepath.Join(pdir, "Makefile"),
			in:       plugin,
			template: nodejsMakefile,
		},
	}
}
