package plugin

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/release"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// ExitError 代表有指定 exit status 的 error
type ExitError struct {
	error
	ExitStatus int
}

// LoadPluginCommands 將 plugin 載入後轉換成 command
func LoadPluginCommands(metadata *release.Metadata) ([]*cobra.Command, error) {
	var commands []*cobra.Command
	if off, _ := strconv.ParseBool(os.Getenv("SL_PLUGINS_OFF")); off {
		return commands, nil
	}
	processParentFlags := func(cmd *cobra.Command, args []string) ([]string, error) {
		k, u := processFlags(args)
		if err := cmd.Parent().ParseFlags(k); err != nil {
			return nil, err
		}
		return u, nil
	}
	found, err := LoadPaths(environment.Settings.PluginDirs())
	if err != nil {
		return commands, fmt.Errorf("failed to load plugins: %s", err)
	}
	for _, plug := range found {
		commands = append(commands, plug.transformToCommand(metadata, processParentFlags))
	}
	return commands, nil
}

func (p *Plugin) transformToCommand(metadata *release.Metadata,
	processParentFlags func(cmd *cobra.Command, args []string) ([]string, error)) *cobra.Command {
	md := p.Metadata
	if md.Usage == "" {
		md.Usage = fmt.Sprintf("The %q plugin", md.Name)
	}
	return &cobra.Command{
		Use:   md.Name,
		Short: md.Usage,
		Long:  md.Description,

		RunE: func(cmd *cobra.Command, args []string) error {
			u, err := processParentFlags(cmd, args)
			if err != nil {
				return err
			}
			if err = p.SetupEnv(metadata); err != nil {
				return err
			}
			command, err := p.Metadata.Exec.GetCommand()
			if err != nil {
				return err
			}
			main, argv, err := p.PrepareCommand(command, u)
			if err != nil {
				return err
			}

			prog := exec.Command(main, argv...)
			prog.Env = os.Environ()
			prog.Stdin = os.Stdin
			prog.Stdout = logrus.StandardLogger().Out
			prog.Stderr = os.Stderr
			if err := prog.Run(); err != nil {
				if eerr, ok := err.(*exec.ExitError); ok {
					os.Stderr.Write(eerr.Stderr)
					status := eerr.Sys().(syscall.WaitStatus)
					return ExitError{
						error:      fmt.Errorf("plugin %q exited with error", md.Name),
						ExitStatus: status.ExitStatus(),
					}
				}
				return err
			}
			return nil
		},
		// This passes all the flags to the subcommand.
		DisableFlagParsing: true,
	}
}

// 將當前的 flag 依照 environment.IsGlobalFlag 分類成 global 及 local flags
func processFlags(args []string) (global []string, local []string) {
	for i := 0; i < len(args); i++ {
		if a := args[i]; environment.IsGlobalFlag(a) {
			global = append(global, a)
		} else {
			local = append(local, a)
		}
	}
	return
}

// LoadPaths 將傳入的 paths 以 filepath.SplitList 切割後, 依序呼叫 LoadAll 載入所有 plugin
func LoadPaths(paths string) ([]*Plugin, error) {
	var found []*Plugin
	// Let's get all UNIXy and allow path separators
	for _, p := range filepath.SplitList(paths) {
		matches, err := LoadAll(p)
		if err != nil {
			return matches, err
		}
		found = append(found, matches...)
	}
	return found, nil
}

// LoadAll 載入 basedir 中的所有子目錄 plugin (子目錄只會收尋一層)
func LoadAll(basedir string) ([]*Plugin, error) {
	var plugins []*Plugin
	// We want basedir/*/plugin.yaml
	scanpath := filepath.Join(basedir, "*", MetadataFileName)
	matches, err := filepath.Glob(scanpath)
	if err != nil {
		return plugins, err
	}

	if matches == nil {
		return plugins, nil
	}

	for _, yaml := range matches {
		dir := filepath.Dir(yaml)
		p, err := LoadDir(dir)
		if err != nil {
			return plugins, err
		}
		plugins = append(plugins, p)
	}
	return plugins, nil
}

// LoadDir 載入 plugin 資料夾
func LoadDir(dirname string) (*Plugin, error) {
	data, err := ioutil.ReadFile(filepath.Join(dirname, MetadataFileName))
	if err != nil {
		return nil, err
	}
	plug := &Plugin{Dir: dirname}
	if err := yaml.Unmarshal(data, &plug.Metadata); err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(filepath.Join(dirname, SourceFileName))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	plug.Source = string(b)
	plug.Mount = filepath.Join(environment.Settings.Home.Mounts(), plug.Metadata.Name)
	return plug, nil
}
