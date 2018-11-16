# slctl

*Slctl*, stands for SoftLeader Control, is a command line interface for running commands against SoftLeader Services

```sh
Slctl against SoftLeader services.

To begin working with slctl, run the 'slctl init' command:

	$ slctl init

It will set up any necessary local configuration.

Common actions from this point include:

Environment:
  $SL_HOME           set an alternative location for slctl files. By default, these are stored in ~/.sl

Usage:
  slctl [command]

Available Commands:
  help        Help about any command
  init        initialize slctl

Flags:
      --debug         enable verbose output
  -h, --help          help for slctl
      --home string   location of your config. Overrides $SL_HOME (default "/Users/matt/.sl")

Use "slctl [command] --help" for more information about a command.
```

## Feature

- [ ] Port to [Homebrew](https://brew.sh/index_zh-tw) for mac users
- [ ] Port to [Chocolatey](https://chocolatey.org/) for windows users
- [ ] Plugins support
