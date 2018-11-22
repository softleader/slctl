package cmd

import (
	"fmt"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

type PluginError struct {
	error
	Code int
}

func loadPlugins(baseCmd *cobra.Command, out io.Writer) {

	if os.Getenv("SL_NO_PLUGINS") == "1" {
		return
	}

	found, err := findPlugins(settings.PluginDirs())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load plugins: %s", err)
		return
	}

	processParent := func(cmd *cobra.Command, args []string) ([]string, error) {
		k, u := manuallyProcessArgs(args)
		if err := cmd.Parent().ParseFlags(k); err != nil {
			return nil, err
		}
		return u, nil
	}

	// Now we create commands for all of these.
	for _, plug := range found {
		plug := plug
		md := plug.Metadata
		if md.Usage == "" {
			md.Usage = fmt.Sprintf("The %q plugin", md.Name)
		}

		c := &cobra.Command{
			Use:   md.Name,
			Short: md.Usage,
			Long:  md.Description,
			RunE: func(cmd *cobra.Command, args []string) error {
				u, err := processParent(cmd, args)
				if err != nil {
					return err
				}

				// Call setupEnv before PrepareCommand because
				// PrepareCommand uses os.ExpandEnv and expects the
				// setupEnv vars.
				if err = plugin.SetupPluginEnv(settings, md.Name, plug.Dir); err != nil {
					return err
				}
				main, argv := plug.PrepareCommand(u)

				prog := exec.Command(main, argv...)
				prog.Env = os.Environ()
				prog.Stdin = os.Stdin
				prog.Stdout = out
				prog.Stderr = os.Stderr
				if err := prog.Run(); err != nil {
					if eerr, ok := err.(*exec.ExitError); ok {
						os.Stderr.Write(eerr.Stderr)
						status := eerr.Sys().(syscall.WaitStatus)
						return PluginError{
							error: fmt.Errorf("plugin %q exited with error", md.Name),
							Code:  status.ExitStatus(),
						}
					}
					return err
				}
				return nil
			},
			// This passes all the flags to the subcommand.
			DisableFlagParsing: true,
		}

		// TODO: Make sure a command with this name does not already exist.
		baseCmd.AddCommand(c)
	}
}

// manuallyProcessArgs processes an arg array, removing special args.
//
// Returns two sets of args: known and unknown (in that order)
func manuallyProcessArgs(args []string) (known []string, unknown []string) {
	for i := 0; i < len(args); i++ {
		switch a := args[i]; a {
		case "--verbose":
			known = append(known, a)
		case "-v":
			known = append(known, a)
		default:
			unknown = append(unknown, a)
		}
	}
	return
}

// findPlugins returns a list of YAML files that describe plugins.
func findPlugins(plugdirs string) ([]*plugin.Plugin, error) {
	found := []*plugin.Plugin{}
	// Let's get all UNIXy and allow path separators
	for _, p := range filepath.SplitList(plugdirs) {
		matches, err := plugin.LoadAll(p)
		if err != nil {
			return matches, err
		}
		found = append(found, matches...)
	}
	return found, nil
}
