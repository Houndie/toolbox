package toolbox

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Add adds a new tool found at packageName to the vendoring system
func Add(packageName string, options ...Option) error {
	return AddVer(packageName, "", options...)
}

// AddVer adds a new tool found at packageName with a specific version to the vendoring system
func AddVer(packageName, version string, options ...Option) error {
	p, err := parseOptions(options...)
	if err != nil {
		return fmt.Errorf("error parsing options: %w", err)
	}

	pkgVer := packageName
	if version != "" {
		pkgVer = packageName + "@" + version
	}

	args := []string{"get"}
	if p.buildFlags != "" {
		args = append(args, strings.Fields(p.buildFlags)...)
	}
	args = append(args, pkgVer)

	goget := exec.Command(p.goBinary, args...)
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

	tools, err := readTools(p.toolsfileName)
	if err != nil {
		return err
	}

	needsUpdate := true
	for _, tool := range tools {
		if tool.Pkg == packageName {
			if tool.BuildFlags == p.buildFlags {
				needsUpdate = false
			}
			break
		}
	}
	if needsUpdate {
		newTool := &tool{
			Pkg:        packageName,
			BuildFlags: p.buildFlags,
		}
		tools = append(tools, newTool)
		if err := writeTools(tools, p.toolsfileName, p.goimportsBinary); err != nil {
			return err
		}
	}

	return nil
}
