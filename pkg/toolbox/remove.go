package toolbox

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/kballard/go-shellquote"
)

// Remove stops tracking the tool from packageName in our vendoring system.
func Remove(packageName string, options ...Option) error {
	p, err := parseOptions(options...)
	if err != nil {
		return fmt.Errorf("error parsing options: %w", err)
	}

	tools, err := readTools(p)
	if err != nil {
		return err
	}

	for i, tool := range tools {
		if tool.Pkg == packageName {
			tools[len(tools)-1], tools[i] = tools[i], tools[len(tools)-1]
			if err := writeTools(tools[:len(tools)-1], p); err != nil {
				return err
			}
			break
		}
	}
	_, dependencyPkg := path.Split(packageName)
	dependencyFile := filepath.Join(p.toolsdirName, dependencyPkg)
	if _, err := os.Stat(dependencyFile); !os.IsNotExist(err) {
		p.logger.Printf("removing file %s", dependencyFile)
		if err := os.Remove(dependencyFile); err != nil {
			return fmt.Errorf("error deleting dependency executable: %w", err)
		}
	} else {
		p.logger.Printf("could not find file %s for removal", dependencyFile)
	}

	gomod := exec.Command(p.goBinary, "mod", "tidy", "-v")
	gomod.Stdout = newLogWriter(p.logger)
	gomod.Stderr = newLogWriter(p.logger)
	p.logger.Printf("calling \"%s\"", shellquote.Join(gomod.Args...))
	if err := gomod.Run(); err != nil {
		return fmt.Errorf("error calling go mod tidy: %w", err)
	}

	return nil
}
