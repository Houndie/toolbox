package toolbox

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

type Tool struct {
	Package    string `json:"package"`
	Version    string `json:"version"`
	BuildFlags string `json:"build_flags"`
}

func List(options ...Option) ([]*Tool, error) {
	p, err := parseOptions(options...)
	if err != nil {
		return nil, fmt.Errorf("error parsing options: %w", err)
	}

	tools, err := readTools(p)
	if err != nil {
		return nil, err
	}

	filename := filepath.Join(p.basedirName, "go.mod")
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening modfile %s: %w", filename, err)
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading in modflie %s: %w", filename, err)
	}
	parseFile, err := modfile.Parse(filename, bytes, nil)
	if err != nil {
		return nil, fmt.Errorf("error parsing modfile %s: %w", filename, err)
	}

	retVals := make([]*Tool, len(tools))
	for i, t := range tools {
		var version string
		for _, m := range parseFile.Require {
			if strings.HasPrefix(t.Pkg, m.Mod.Path) {
				version = m.Mod.Version
				break
			}
		}
		if version == "" {
			return nil, fmt.Errorf("no version for package %s found", t.Pkg)
		}

		retVals[i] = &Tool{
			Package:    t.Pkg,
			Version:    version,
			BuildFlags: t.BuildFlags,
		}

	}
	return retVals, nil
}
