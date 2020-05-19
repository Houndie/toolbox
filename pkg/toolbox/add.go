package toolbox

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Add adds a new tool found at packageName to the vendoring system
func Add(packageName string, options ...Option) error {
	return AddVer(packageName, "", options...)
}

// AddVer adds a new tool found at packageName with a specific version to the vendoring system
func AddVer(packageName, version string, options ...Option) error {
	p := parseOptions(options...)

	tools, err := readTools(p.toolsfileName)
	if err != nil {
		return err
	}

	found := false
	for _, tool := range tools {
		if tool == packageName {
			found = true
			break
		}
	}
	if !found {
		tools = append(tools, packageName)
		if err := writeTools(tools, p.toolsfileName, p.goimportsBinary); err != nil {
			return err
		}
	}

	if version != "" {
		packageName = packageName + "@" + version
	}

	goget := exec.Command(p.goBinary, "get", packageName)
	absToolsdir, err := filepath.Abs(p.toolsdirName)
	if err != nil {
		return fmt.Errorf("error finding absolute path to toolsdir %s: %w", p.toolsdirName, err)
	}
	goget.Env = append(os.Environ(), "GOBIN="+absToolsdir)
	if _, err := goget.Output(); err != nil {
		eerr := &exec.ExitError{}
		if !errors.As(err, &eerr) {
			return fmt.Errorf("error calling go get: %w", err)
		}
		return fmt.Errorf("error calling go get: %s: %w", string(eerr.Stderr), err)
	}
	return nil
}
