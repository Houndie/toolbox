package toolbox

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kballard/go-shellquote"
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

	args := []string{"get", "-v"}
	if p.buildFlags != "" {
		split, err := shellquote.Split(p.buildFlags)
		if err != nil {
			return fmt.Errorf("error splitting args: %w", err)
		}
		args = append(args, split...)
	}
	args = append(args, pkgVer)

	goget := exec.Command(p.goBinary, args...)
	absToolsdir, err := filepath.Abs(p.toolsdirName)
	if err != nil {
		return fmt.Errorf("error finding absolute path to toolsdir %s: %w", p.toolsdirName, err)
	}
	goget.Env = append(os.Environ(), "GOBIN="+absToolsdir)
	p.logger.Printf("calling \"%s\", with GOBIN=%s", strings.Join(goget.Args, " "), absToolsdir)
	goget.Stdout = newLogWriter(p.logger)
	goget.Stderr = newLogWriter(p.logger)
	if err := goget.Run(); err != nil {
		return fmt.Errorf("error calling go get: %w", err)
	}

	tools, err := readTools(p)
	if err != nil {
		return err
	}

	needsUpdate := true
	found := false
	for _, tool := range tools {
		if tool.Pkg == packageName {
			if tool.BuildFlags == p.buildFlags {
				needsUpdate = false
			} else {
				tool.BuildFlags = p.buildFlags
			}
			found = true
			break
		}
	}
	if !found {
		tools = append(tools, &tool{
			Pkg:        packageName,
			BuildFlags: p.buildFlags,
		})
	}
	if needsUpdate {
		if err := writeTools(tools, p); err != nil {
			return err
		}
	}

	return nil
}
