package main

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const toolsfile = "tools.go"

var toolsTemplate = template.Must(template.New("tools_template").Parse(`// +build tools

// This file is generated and managed by toolbox.  Manually edit at your own peril.
package toolbox

import (
	{{range .}}
	_ "{{ . }}" {{end}}
)`))

func readTools() ([]string, error) {
	if _, err := os.Stat(toolsfile); os.IsNotExist(err) {
		return []string{}, nil
	}

	file, err := parser.ParseFile(token.NewFileSet(), toolsfile, nil, parser.ImportsOnly)
	if err != nil {
		return nil, fmt.Errorf("error parsing tools file %s: %w", toolsfile, err)
	}

	tools := make([]string, len(file.Imports))
	for i, imp := range file.Imports {
		tools[i] = strings.TrimSuffix(strings.TrimPrefix(imp.Path.Value, "\""), "\"")
	}
	return tools, nil
}

func writeTools(tools []string) error {
	file, err := os.OpenFile(toolsfile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("error opening tools file %s: %w", toolsfile, err)
	}
	defer file.Close()

	if err := toolsTemplate.Execute(file, tools); err != nil {
		return fmt.Errorf("error writing data to toolsfile %s: %w", toolsfile, err)
	}

	if _, err := exec.Command("goimports", "-w", toolsfile).Output(); err != nil {
		eerr := &exec.ExitError{}
		if !errors.As(err, &eerr) {
			return fmt.Errorf("error calling goimports: %w", err)
		}
		return fmt.Errorf("error calling goimports: %s: %w", string(eerr.Stderr), err)
	}

	return nil
}
