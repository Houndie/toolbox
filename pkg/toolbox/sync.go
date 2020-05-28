package toolbox

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		if p.buildFlags != "" {
			args = append(args, strings.Fields(p.buildFlags)...)
		}
		args = append(args, t.Pkg)
		goinstall := exec.Command(p.goBinary, "install", t.Pkg)
		goinstall.Stdout = newLogWriter(p.logger)
		goinstall.Stderr = newLogWriter(p.logger)
		absToolsdir, err := filepath.Abs(p.toolsdirName)
		if err != nil {
			return fmt.Errorf("error finding absolute path to toolsdir %s: %w", p.toolsdirName, err)
		}
		goinstall.Env = append(os.Environ(), "GOBIN="+absToolsdir)
		p.logger.Printf("running \"%s\", with GOBIN=%s", strings.Join(goinstall.Args, " "), absToolsdir)
		if err := goinstall.Run(); err != nil {
			return fmt.Errorf("error calling go install: %w", err)
		}
	}

	return nil
}
