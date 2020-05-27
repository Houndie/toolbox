package toolbox

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Do runs the given command, using a vendored tool if applicable
func Do(command ...string) error {
	return DoOpts(command)
}

// Do runs the given command with the given options, using a vendored tool if applicable
func DoOpts(command []string, options ...Option) error {
	if len(command) < 1 {
		return nil
	}

	cmd, err := CommandOpts(command[0], command[1:], options...)
	if err != nil {
		return err
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// We don't really care about the error from this, stderr goes to screen.
	_ = cmd.Run()
	return nil
}

// Command creates an exec.Cmd, modifying the environment and run binary name to reflect the vendored tools.
func Command(command string, args ...string) (*exec.Cmd, error) {
	return CommandOpts(command, args)
}

// CommandOpts creates an exec.Cmd, modifying the environment and run binary name to reflect the vendored tools, using the given options.
func CommandOpts(command string, args []string, options ...Option) (*exec.Cmd, error) {
	p := parseOptions(options...)

	absToolsdir, err := filepath.Abs(p.toolsdirName)
	if err != nil {
		return nil, fmt.Errorf("error finding absolute path to toolsdir %s: %w", p.toolsdirName, err)
	}

	doCommand := command
	if !strings.Contains(command, string(filepath.Separator)) {
		potentialCommand := filepath.Join(absToolsdir, command)
		if _, err := os.Stat(potentialCommand); err == nil {
			doCommand = potentialCommand
		}
	}

	cmd := exec.Command(doCommand, args...)

	// GOBIN is set so that tools that install other tools still prefer the tool directory
	cmd.Env = append(os.Environ(),
		"GOBIN="+absToolsdir,
		"PATH="+absToolsdir+string(filepath.ListSeparator)+os.Getenv("PATH"),
	)

	return cmd, nil
}
