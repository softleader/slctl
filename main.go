package main

import (
	"github.com/softleader/slctl/cmd"
	"os"
)

func main() {
	command := cmd.NewRootCmd(os.Args[1:])
	if err := command.Execute(); err != nil {
		switch e := err.(type) {
		case cmd.PluginError:
			os.Exit(e.Code)
		default:
			os.Exit(1)
		}
	}
}
