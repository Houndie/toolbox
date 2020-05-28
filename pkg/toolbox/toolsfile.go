package toolbox

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/kballard/go-shellquote"
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

func readTools(p *parsedOptions) ([]*tool, error) {
	if _, err := os.Stat(p.toolsfileName); os.IsNotExist(err) {
		return nil, nil
	}

	p.logger.Printf("parsing toolsfile %s", p.toolsfileName)
	file, err := parser.ParseFile(token.NewFileSet(), p.toolsfileName, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing tools file %s: %w", p.toolsfileName, err)
	}

	tools := make([]*tool, len(file.Imports))
	for i, imp := range file.Imports {
		tools[i] = &tool{}
		if imp.Comment.Text() != "" {
			if err := json.Unmarshal([]byte(imp.Comment.Text()), &tools[i]); err != nil {
				return nil, fmt.Errorf("error parsing tool comment as json: %w", err)
			}
		}
		tools[i].Pkg = strings.TrimSuffix(strings.TrimPrefix(imp.Path.Value, "\""), "\"")
	}
	return tools, nil
}

func writeTools(tools []*tool, p *parsedOptions) error {
	file, err := os.OpenFile(p.toolsfileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("error opening tools file %s: %w", p.toolsfileName, err)
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

	p.logger.Printf("writing toolsfile %s", p.toolsfileName)
	if err := toolsTemplate.Execute(file, toolTemplates); err != nil {
		return fmt.Errorf("error writing data to toolsfile %s: %w", p.toolsfileName, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("error closing toolsfile %s: %w", p.toolsfileName, err)
	}

	if _, err := exec.LookPath(p.goimportsBinary); err == nil {
		goimports := exec.Command(p.goimportsBinary, "-v", "-w", p.toolsfileName)
		goimports.Stdout = newLogWriter(p.logger)
		goimports.Stderr = newLogWriter(p.logger)
		p.logger.Printf("calling \"%s\"", shellquote.Join(goimports.Args...))
		if err := goimports.Run(); err != nil {
			return fmt.Errorf("error calling goimports: %w", err)
		}
	}

	return nil
}
