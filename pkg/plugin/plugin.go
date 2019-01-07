package plugin

import (
	"os"
	"strings"
)

const MetadataFileName = "metadata.yaml"
const SourceFileName = ".source"

type Metadata struct {
	Name              string   `json:"name"`
	Version           string   `json:"version"`
	Usage             string   `json:"usage"`
	Description       string   `json:"description"`
	Exec              Commands `json:"exec"`
	Hook              Commands `json:"hook"`
	IgnoreGlobalFlags bool     `json:"ignoreGlobalFlags"`
	GitHub            GitHub   `json:"github"`
}

type Plugin struct {
	Metadata *Metadata
	Dir      string
	Source   string // 在安裝 plugin 時的 source, 只有非本機的 source 才會紀錄, 為了方便之後做 github plugin 的 upgrade
	Mount    string
}

func (p *Plugin) FromGitHub() bool {
	return strings.HasPrefix(p.Source, "github.com/")
}

// PrepareCommand takes a Plugin.Command and prepares it for execution.
//
// It merges extraArgs into any arguments supplied in the plugin. It
// returns the name of the command and an args array.
//
// The result is suitable to pass to exec.Command.
func (p *Plugin) PrepareCommand(command string, extraArgs []string) (main string, argv []string, err error) {
	parts := strings.Split(os.ExpandEnv(command), " ")
	main = parts[0]
	if len(parts) > 1 {
		argv = parts[1:]
	}
	if !p.Metadata.IgnoreGlobalFlags && extraArgs != nil {
		argv = append(argv, extraArgs...)
	}
	return
}

