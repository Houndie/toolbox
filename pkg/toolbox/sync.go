package toolbox

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Sync interates through all tools that we're vendoring, and ensures that all of them are installed, and at the correct version.
func Sync(options ...Option) error {
	p := parseOptions(options...)
	tools, err := readTools(p.toolsfileName)
	if err != nil {
		return err
	}

	for _, tool := range tools {
		goinstall := exec.Command(p.goBinary, "install", tool)
		absToolsdir, err := filepath.Abs(p.toolsdirName)
		if err != nil {
			return fmt.Errorf("error finding absolute path to toolsdir %s: %w", p.toolsdirName, err)
		}
		goinstall.Env = append(os.Environ(), "GOBIN="+absToolsdir)
		if _, err := goinstall.Output(); err != nil {
			eerr := &exec.ExitError{}
			if !errors.As(err, &eerr) {
				return fmt.Errorf("error calling go install: %w", err)
			}
			return fmt.Errorf("error calling go install: %s: %w", string(eerr.Stderr), err)
		}
	}

	return nil
}
