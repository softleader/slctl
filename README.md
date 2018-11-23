# slctl

Slctl is a command line interface for running commands against SoftLeader services.

## Install

You can also use [Homebrew](https://brew.sh/index_zh-tw) (on macOS):

```sh
brew install softleader/tap/slctl
```

Or install using [Chocolatey](https://chocolatey.org/) (on Windows):

```sh
TODO
```

Or manually downlaod from [releases page](https://github.com/softleader/slctl/releases).

## Usage

```sh
To begin working with slctl, run the 'slctl init' command:

	$ slctl init

It will set up any necessary local configuration.

Common actions from this point include:

Environment:
  $SL_HOME           set an alternative location for slctl files. By default, these are stored in ~/.sl
  $SL_NO_PLUGINS     disable plugins. Set $SL_NO_PLUGINS=true to disable plugins.
  $SL_OFFLINE   	 work offline. Set $SL_OFFLINE=true to work offline.

Usage:
  slctl [command]

Available Commands:
  help        Help about any command
  home        displays the location of SL_HOME
  init        initialize slctl
  plugin      add, list, remove, or create plugins
  version     print slctl version.

Flags:
  -h, --help          help for slctl
      --home string   location of your config. Overrides $SL_HOME (default "~/.sl")
      --offline       work offline
  -v, --verbose       enable verbose output

Use "slctl [command] --help" for more information about a command.
```

## Plugins

- [foo](#) - The foo plugin
- [bar](#) - The bar plugin

> TODO: Plugin 清單

### Writing Custom Plugins

*Slctl* 支援任何語言的 Plugin, 請參考 [Plugins Guide](https://github.com/softleader/slctl/wiki/Plugins-Guide)

