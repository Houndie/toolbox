package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "toolbox",
	Short: "Toolbox is a tool vendoring helper",
	Long:  "Toolbox sits on top of go's powerful module engine, and leverages it to vendor your executables",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
