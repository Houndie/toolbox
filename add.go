package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addCommand = &cobra.Command{
	Use:   "add <dependency> [version]",
	Short: "Add a new dependency",
	Long:  "Adds dependency to the list of dependencies managed by toolbox.  If a version is provided, adds that version as well.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dependency := args[0]

		tools, err := readTools()
		if err != nil {
			return err
		}

		found := false
		for _, tool := range tools {
			if tool == dependency {
				found = true
				break
			}
		}
		if !found {
			tools = append(tools, dependency)
			if err := writeTools(tools); err != nil {
				return err
			}
		}

		if len(args) > 1 {
			dependency = dependency + "@" + args[1]
		}

		goget := exec.Command(viper.GetString(goFlag), "get", dependency)
		toolsdir, err := toolsDir()
		if err != nil {
			return err
		}
		goget.Env = append(os.Environ(), "GOBIN="+toolsdir)
		if _, err := goget.Output(); err != nil {
			eerr := &exec.ExitError{}
			if !errors.As(err, &eerr) {
				return fmt.Errorf("error calling go get: %w", err)
			}
			return fmt.Errorf("error calling go get: %s: %w", string(eerr.Stderr), err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCommand)
}
