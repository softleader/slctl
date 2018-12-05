package plugin

import (
	"fmt"
	"runtime"
)

type ErrNoCommandFound struct {
	s string
}

func (e *ErrNoCommandFound) Error() string {
	return e.s
}

type Commands struct {
	Command  string     `json:"command"`
	Platform []Platform `json:"platform"`
}

type Platform struct {
	Os      string `json:"os"`
	Arch    string `json:"arch"`
	Command string `json:"command"`
}

func (c *Commands) GetCommand() (command string, err error) {
	command = c.Command
	for _, p := range c.Platform {
		if p.Os == runtime.GOOS && p.Arch == runtime.GOARCH {
			command = p.Command
			return
		}
		if p.Os == runtime.GOOS {
			command = p.Command
		}
	}
	if command == "" {
		err = &ErrNoCommandFound{
			s: fmt.Sprintf("no suitable command found for %s/%s", runtime.GOOS, runtime.GOARCH),
		}
	}
	return
}
