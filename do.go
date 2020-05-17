package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var doCommand = &cobra.Command{
	Use:   "do <command>",
	Short: "runs a command using the vendored version of tools",
	Long:  "do edits the path to reflect the tool vendor directly, and runs the given command",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return nil
		}

		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error fetching current working directory: %w", err)
		}

		if err := os.Setenv("PATH", filepath.Join(cwd, "_tools")+string(os.PathListSeparator)+os.Getenv("PATH")); err != nil {
			return fmt.Errorf("error setting path: %w", err)
		}

		docmd := exec.Command(args[0], args[1:]...)

		docmd.Stdin = os.Stdin
		docmd.Stdout = os.Stdout
		docmd.Stderr = os.Stderr

		// GOBIN is set so that tools that install other tools still prefer the tool directory
		docmd.Env = append(os.Environ(),
			"GOBIN="+filepath.Join(cwd, "_tools"),
		)

		// We don't really care about the error from this, stderr goes to screen.
		_ = docmd.Run()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doCommand)
}
