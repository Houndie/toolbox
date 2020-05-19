package main

import (
	"github.com/Houndie/toolbox/pkg/toolbox"
	"github.com/spf13/cobra"
)

var doCommand = &cobra.Command{
	Use:   "do <command>",
	Short: "Run a command using the vendored version of tools",
	Long:  "Edits the PATH to reflect the tool vendor directly, and runs the given command.",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := toolbox.DoOpts(args, makeOptions()...)
		if err == nil {
			return nil
		}
		return err
	},
}

var addCommand = &cobra.Command{
	Use:   "add <dependency> [version]",
	Short: "Add a new dependency",
	Long:  "Adds dependency to the list of dependencies managed by toolbox.  If a version is provided, adds that version as well.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return toolbox.AddVer(args[0], args[1], makeOptions()...)
		}
		return toolbox.Add(args[0], makeOptions()...)
	},
}

var removeCommand = &cobra.Command{
	Use:   "remove <dependency>",
	Short: "Remove a dependency",
	Long:  "Removes a dependency, and attempts to remove the executable of the same name as well.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return toolbox.Remove(args[0], makeOptions()...)
	},
}

var syncCommand = &cobra.Command{
	Use:   "sync",
	Short: "Make sure all dependencies are at the correct version",
	Long:  "Uses go install to install all of our dependencies.  Installs from module cache if they are found, from the internet if not.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return toolbox.Sync(makeOptions()...)
	},
}

func init() {
	rootCmd.AddCommand(doCommand)
	rootCmd.AddCommand(addCommand)
	rootCmd.AddCommand(removeCommand)
	rootCmd.AddCommand(syncCommand)
}
