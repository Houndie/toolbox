package toolbox

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// The default option values
const (
	defaultGo        = "go"
	defaultGoimports = "goimports"
)

func defaultBasedir(goCommand string) (string, error) {
	out, err := exec.Command(goCommand, "env", "GOMOD").Output()
	if err != nil {
		return "", fmt.Errorf("error finding module root: %w", err)
	}
	if string(out) == "" {
		return "", fmt.Errorf("no go module found, please initialize with \"go mod init\"")
	}
	return filepath.Dir(string(out)), nil
}

func defaultToolsfile(basedir string) string {
	return filepath.Join(basedir, "tools.go")
}

func defaultToolsdir(basedir string) string {
	return filepath.Join(basedir, "_tools")
}

type parsedOptions struct {
	goBinary        string
	goimportsBinary string
	toolsfileName   string
	toolsdirName    string
	basedirName     string
}

// Option is an optional modifier to toolbox's default behavior
type Option interface {
	apply(*parsedOptions) *parsedOptions
}

type goOption struct {
	goBinary string
}

func (o *goOption) apply(p *parsedOptions) *parsedOptions {
	p.goBinary = o.goBinary
	return p
}

// GoOption changes the default name/path of the "go" binary.
func GoOption(goBinaryName string) Option {
	return &goOption{goBinary: goBinaryName}
}

type goimportsOption struct {
	goimportsBinary string
}

func (o *goimportsOption) apply(p *parsedOptions) *parsedOptions {
	p.goimportsBinary = o.goimportsBinary
	return p
}

// GoimportsOption changes the default name/path of the "goimports" binary.
func GoimportsOption(goimportsBinaryName string) Option {
	return &goimportsOption{goimportsBinary: goimportsBinaryName}
}

type toolsfileOption struct {
	toolsfileName string
}

func (o *toolsfileOption) apply(p *parsedOptions) *parsedOptions {
	p.toolsfileName = o.toolsfileName
	return p
}

// ToolfileOption changes the default name/path of the file used to manage tool dependencies
func ToolsfileOption(toolsfileName string) Option {
	return &toolsfileOption{toolsfileName: toolsfileName}
}

type toolsdirOption struct {
	toolsdirName string
}

func (o *toolsdirOption) apply(p *parsedOptions) *parsedOptions {
	p.toolsdirName = o.toolsdirName
	return p
}

// TooldirOption changes the default name/path of the directory used to vendor tools
func ToolsdirOption(toolsdirName string) Option {
	return &toolsdirOption{toolsdirName: toolsdirName}
}

type basedirOption struct {
	basedirName string
}

func (o *basedirOption) apply(p *parsedOptions) *parsedOptions {
	p.basedirName = o.basedirName
	return p
}

// TooldirOption changes the default name/path of the directory used to vendor tools
func BasedirOption(basedirName string) Option {
	return &basedirOption{basedirName: basedirName}
}

func parseOptions(options ...Option) (*parsedOptions, error) {
	p := &parsedOptions{}
	for _, option := range options {
		p = option.apply(p)
	}

	if p.goBinary == "" {
		p.goBinary = defaultGo
	}
	if p.goimportsBinary == "" {
		p.goimportsBinary = defaultGoimports
	}
	if p.basedirName == "" {
		var err error
		p.basedirName, err = defaultBasedir(p.goBinary)
		if err != nil {
			return nil, fmt.Errorf("error determining default base directory: %w", err)
		}
	}
	if p.toolsdirName == "" {
		p.toolsdirName = defaultToolsdir(p.basedirName)
	}
	if p.toolsfileName == "" {
		p.toolsfileName = defaultToolsfile(p.basedirName)
	}
	return p, nil
}
