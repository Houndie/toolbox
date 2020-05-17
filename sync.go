package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var syncCommand = &cobra.Command{
	Use:   "sync",
	Short: "makes sure all dependencies are at the correct version",
	Long:  "uses go install to install all of our dependencies.  Installs from module cache if they are found, from the internet if not",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		tools, err := readTools()
		if err != nil {
			return err
		}

		for _, tool := range tools {
			goinstall := exec.Command(viper.GetString(goFlag), "install", tool)
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("error fetching current working directory: %w", err)
			}
			goinstall.Env = append(os.Environ(), "GOBIN="+filepath.Join(cwd, "_tools"))
			if _, err := goinstall.Output(); err != nil {
				eerr := &exec.ExitError{}
				if !errors.As(err, &eerr) {
					return fmt.Errorf("error calling go install: %w", err)
				}
				return fmt.Errorf("error calling go install: %s: %w", string(eerr.Stderr), err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCommand)
}
