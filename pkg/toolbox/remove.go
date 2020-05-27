package toolbox

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

// Remove stops tracking the tool from packageName in our vendoring system.
func Remove(packageName string, options ...Option) error {
	p, err := parseOptions(options...)
	if err != nil {
		return fmt.Errorf("error parsing options: %w", err)
	}

	tools, err := readTools(p.toolsfileName)
	if err != nil {
		return err
	}

	for i, tool := range tools {
		if tool == packageName {
			tools[len(tools)-1], tools[i] = tools[i], tools[len(tools)-1]
			if err := writeTools(tools[:len(tools)-1], p.toolsfileName, p.goimportsBinary); err != nil {
				return err
			}
			break
		}
	}
	_, dependencyPkg := path.Split(packageName)
	dependencyFile := filepath.Join(p.toolsdirName, dependencyPkg)
	if _, err := os.Stat(dependencyFile); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(dependencyFile); err != nil {
		return fmt.Errorf("error deleting dependency executable: %w", err)
	}

	if _, err := exec.Command(p.goBinary, "mod", "tidy").Output(); err != nil {
		eerr := &exec.ExitError{}
		if !errors.As(err, &eerr) {
			return fmt.Errorf("error calling go mod tidy: %w", err)
		}
		return fmt.Errorf("error calling go mod tidy: %s: %w", string(eerr.Stderr), err)
	}

	return nil
}
