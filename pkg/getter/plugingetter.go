package getter

import (
	"bytes"
	"fmt"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/plugin"
	"os"
	"os/exec"
	"path/filepath"

)

// collectPlugins scans for getter plugins.
// This will load plugins according to the environment.
func collectPlugins(settings environment.EnvSettings) (Providers, error) {
	plugins, err := plugin.FindPlugins(settings.PluginDirs())
	if err != nil {
		return nil, err
	}
	var result Providers
	for _, plugin := range plugins {
		for _, downloader := range plugin.Metadata.Downloaders {
			result = append(result, Provider{
				Schemes: downloader.Protocols,
				New: newPluginGetter(
					downloader.Command,
					settings,
					plugin.Metadata.Name,
					plugin.Dir,
				),
			})
		}
	}
	return result, nil
}

// pluginGetter is a generic type to invoke custom downloaders,
// implemented in plugins.
type pluginGetter struct {
	command                   string
	certFile, keyFile, cAFile string
	settings                  environment.EnvSettings
	name                      string
	base                      string
}

// Get runs downloader plugin command
func (p *pluginGetter) Get(href string) (*bytes.Buffer, error) {
	argv := []string{p.certFile, p.keyFile, p.cAFile, href}
	prog := exec.Command(filepath.Join(p.base, p.command), argv...)
	plugin.SetupPluginEnv(p.settings, p.name, p.base)
	prog.Env = os.Environ()
	buf := bytes.NewBuffer(nil)
	prog.Stdout = buf
	prog.Stderr = os.Stderr
	prog.Stdin = os.Stdin
	if err := prog.Run(); err != nil {
		if eerr, ok := err.(*exec.ExitError); ok {
			os.Stderr.Write(eerr.Stderr)
			return nil, fmt.Errorf("plugin %q exited with error", p.command)
		}
		return nil, err
	}
	return buf, nil
}

// newPluginGetter constructs a valid plugin getter
func newPluginGetter(command string, settings environment.EnvSettings, name, base string) Constructor {
	return func(URL, CertFile, KeyFile, CAFile string) (Getter, error) {
		result := &pluginGetter{
			command:  command,
			certFile: CertFile,
			keyFile:  KeyFile,
			cAFile:   CAFile,
			settings: settings,
			name:     name,
			base:     base,
		}
		return result, nil
	}
}
