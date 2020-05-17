package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var removeCommand = &cobra.Command{
	Use:   "remove <dependency>",
	Short: "removes a dependency",
	Long:  "removes a dependency, and attempts to remove the executable of the same name as well",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dependency := args[0]

		tools, err := readTools()
		if err != nil {
			return err
		}

		for i, tool := range tools {
			if tool == dependency {
				tools[len(tools)-1], tools[i] = tools[i], tools[len(tools)-1]
				if err := writeTools(tools[:len(tools)-1]); err != nil {
					return err
				}
				break
			}
		}
		_, dependencyPkg := path.Split(dependency)
		dependencyFile := filepath.Join("_tools", dependencyPkg)
		if _, err := os.Stat(dependencyFile); os.IsNotExist(err) {
			return nil
		}

		if err := os.Remove(dependencyFile); err != nil {
			return fmt.Errorf("error deleting dependency executable: %w", err)
		}

		if _, err := exec.Command(viper.GetString(goFlag), "mod", "tidy").Output(); err != nil {
			eerr := &exec.ExitError{}
			if !errors.As(err, &eerr) {
				return fmt.Errorf("error calling go mod tidy: %w", err)
			}
			return fmt.Errorf("error calling go mod tidy: %s: %w", string(eerr.Stderr), err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCommand)
}
