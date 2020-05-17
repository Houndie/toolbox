package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

const goFlag = "go"

func init() {
	rootCmd.PersistentFlags().String(goFlag, "go", "the \"go\" executable to use")
	viper.BindPFlag(goFlag, rootCmd.PersistentFlags().Lookup(goFlag))

	cobra.OnInitialize(func() {
		viper.AddConfigPath(".")
		viper.SetConfigName(".toolbox")
		viper.AutomaticEnv()

		_ = viper.ReadInConfig()
	})
}
