// Package toolbox is a library for automating tool dependency management.  It is designed as the logical portion of the toolbox binary.
//
// Toolbox uses a file (default "tools.go") to manage its dependency list.  This is should always be a .go file, as version information is tracked via go modules, and thus tool dependencies need to be kept in a file that the go module system can scan.  The tools themselves are kept in a folder (default "_tools").
//
// Toolbox doesn't do anything massively unique, it simply automates the fiddly bits...it implements the best practices for tool automation found here: https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md
package toolbox
