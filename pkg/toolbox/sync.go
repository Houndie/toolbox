package toolbox

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kballard/go-shellquote"
)

// Sync interates through all tools that we're vendoring, and ensures that all of them are installed, and at the correct version.
func Sync(options ...Option) error {
	p, err := parseOptions(options...)
	if err != nil {
		return fmt.Errorf("error parsing options: %w", err)
	}
	tools, err := readTools(p)
	if err != nil {
		return err
	}

	for _, t := range tools {
		args := []string{"install", "-v"}
		if t.BuildFlags != "" {
			split, err := shellquote.Split(t.BuildFlags)
			if err != nil {
				return fmt.Errorf("error splitting args: %w", err)
			}
			args = append(args, split...)
		}
		args = append(args, t.Pkg)
		goinstall := exec.Command(p.goBinary, args...)
		goinstall.Stdout = newLogWriter(p.logger)
		goinstall.Stderr = newLogWriter(p.logger)
		absToolsdir, err := filepath.Abs(p.toolsdirName)
		if err != nil {
			return fmt.Errorf("error finding absolute path to toolsdir %s: %w", p.toolsdirName, err)
		}
		goinstall.Env = append(os.Environ(), "GOBIN="+absToolsdir)
		p.logger.Printf("running \"%s\", with GOBIN=%s", shellquote.Join(goinstall.Args...), absToolsdir)
		if err := goinstall.Run(); err != nil {
			return fmt.Errorf("error calling go install: %w", err)
		}
	}

	return nil
}
