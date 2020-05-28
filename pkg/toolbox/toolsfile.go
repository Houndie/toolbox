package toolbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

var toolsTemplate = template.Must(template.New("tools_template").Parse(`// +build tools

// This file is generated and managed by toolbox.  Manually edit at your own peril.
package toolbox

import (
	{{range .}}
	_ "{{ .Pkg }}" {{if ne .Comment "{}" }} //{{ .Comment }} {{end}}{{end}}
)`))

type tool struct {
	Pkg        string `json:"-"`
	BuildFlags string `json:"build_flags,omitempty"`
}

type toolTemplate struct {
	Pkg     string
	Comment string
}

func readTools(toolsfile string) ([]*tool, error) {
	if _, err := os.Stat(toolsfile); os.IsNotExist(err) {
		return nil, nil
	}

	file, err := parser.ParseFile(token.NewFileSet(), toolsfile, nil, parser.ImportsOnly)
	if err != nil {
		return nil, fmt.Errorf("error parsing tools file %s: %w", toolsfile, err)
	}

	tools := make([]*tool, len(file.Imports))
	for i, imp := range file.Imports {
		tools[i] = &tool{}
		if imp.Comment.Text() != "" {
			json.Unmarshal([]byte(imp.Comment.Text()), &tools[i])
		}
		tools[i].Pkg = strings.TrimSuffix(strings.TrimPrefix(imp.Path.Value, "\""), "\"")
	}
	return tools, nil
}

func writeTools(tools []*tool, toolsfile string, goimports string) error {
	file, err := os.OpenFile(toolsfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("error opening tools file %s: %w", toolsfile, err)
	}
	defer file.Close()

	toolTemplates := make([]*toolTemplate, len(tools))
	for i, t := range tools {
		j, err := json.Marshal(t)
		if err != nil {
			return fmt.Errorf("error generating toolsfile comment: %w", err)
		}
		toolTemplates[i] = &toolTemplate{
			Pkg:     t.Pkg,
			Comment: string(j),
		}
	}

	if err := toolsTemplate.Execute(file, toolTemplates); err != nil {
		return fmt.Errorf("error writing data to toolsfile %s: %w", toolsfile, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("error closing toolsfile %s: %w", toolsfile, err)
	}

	if _, err := exec.LookPath(goimports); err == nil {
		if _, err := exec.Command(goimports, "-w", toolsfile).Output(); err != nil {
			eerr := &exec.ExitError{}
			if !errors.As(err, &eerr) {
				return fmt.Errorf("error calling goimports: %w", err)
			}
			return fmt.Errorf("error calling goimports: %s: %w", string(eerr.Stderr), err)
		}
	}

	return nil
}
