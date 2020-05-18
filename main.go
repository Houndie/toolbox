package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:     "toolbox",
	Short:   "Toolbox is a tool vendoring helper",
	Long:    "Toolbox sits on top of go's powerful module engine, and leverages it to vendor your executables",
	Version: "v0.1",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

const goFlag = "go"
const goimportsFlag = "goimports"
const toolsfileFlag = "tools_file"
const toolsdirFlag = "tools_directory"
const configfileFlag = "config_file"

func init() {
	rootCmd.PersistentFlags().String(goFlag, "go", "the \"go\" executable to use")
	viper.BindPFlag(goFlag, rootCmd.PersistentFlags().Lookup(goFlag))

	rootCmd.PersistentFlags().String(goimportsFlag, "goimports", "the \"goimports\" executable to use")
	viper.BindPFlag(goimportsFlag, rootCmd.PersistentFlags().Lookup(goimportsFlag))

	rootCmd.PersistentFlags().String(toolsfileFlag, "tools.go", "the file in which to store tool data. This should end in a \".go\" extenstion so that go's module system picks it up.")
	viper.BindPFlag(toolsfileFlag, rootCmd.PersistentFlags().Lookup(toolsfileFlag))

	rootCmd.PersistentFlags().String(toolsdirFlag, "_tools", "the directory where tool binaries are stored")
	viper.BindPFlag(toolsdirFlag, rootCmd.PersistentFlags().Lookup(toolsdirFlag))

	cfgFile := ""
	rootCmd.PersistentFlags().StringVar(&cfgFile, configfileFlag, "", "the location of a config file to load. By default, looks for \".toolbox.ini\", \".toolbox.json\", \".toolbox.yaml\", or \".toolbox.toml\"")

	cobra.OnInitialize(func() {
		if cfgFile == "" {
			viper.AddConfigPath(".")
			viper.SetConfigName(".toolbox")
		} else {
			viper.SetConfigFile(cfgFile)
		}

		viper.SetEnvPrefix("TOOLBOX")
		viper.AutomaticEnv()

		_ = viper.ReadInConfig()
	})
}

func toolsDir() (string, error) {
	toolsdir := viper.GetString(toolsdirFlag)
	if !filepath.IsAbs(toolsdir) {
		return toolsdir, nil
	}

	toolsdir, err := filepath.Abs(toolsdir)
	if err != nil {
		return "", fmt.Errorf("error making absolute tools dir: %w", err)
	}
	return toolsdir, nil
}
