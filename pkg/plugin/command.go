package plugin

import (
	"fmt"
	"runtime"
)

// ErrNoCommandFound 代表找不到此 plugin 在當前環境中可以執行的命令
type ErrNoCommandFound struct {
	s string
}

func (e *ErrNoCommandFound) Error() string {
	return e.s
}

// Commands 封裝了執行 plugin 的執行命令
type Commands struct {
	Command  string     `json:"command"`
	Platform []Platform `json:"platform"`
}

// Platform 封裝了執行 plugin 執行命令的平台資訊
type Platform struct {
	Os      string `json:"os"`
	Arch    string `json:"arch"`
	Command string `json:"command"`
}

// GetCommand 取得符合當前系統環境的執行命令
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
