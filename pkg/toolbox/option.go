package toolbox

// The default option values
const (
	DefaultGo        = "go"
	DefaultGoimports = "goimports"
	DefaultToolsfile = "tools.go"
	DefaultToolsdir  = "_tools"
)

type parsedOptions struct {
	goBinary        string
	goimportsBinary string
	toolsfileName   string
	toolsdirName    string
}

func defaultParsedOptions() *parsedOptions {
	return &parsedOptions{
		goBinary:        DefaultGo,
		goimportsBinary: DefaultGoimports,
		toolsfileName:   DefaultToolsfile,
		toolsdirName:    DefaultToolsdir,
	}
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

func parseOptions(options ...Option) *parsedOptions {
	p := defaultParsedOptions()
	for _, option := range options {
		p = option.apply(p)
	}
	return p
}
