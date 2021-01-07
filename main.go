package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Houndie/toolbox/pkg/toolbox"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:     "toolbox",
	Short:   "Toolbox is a tool vendoring helper",
	Long:    "Toolbox sits on top of go's powerful module engine, and leverages it to vendor your executables",
	Version: "v0.2",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		execErr := &exec.ExitError{}
		if errors.As(err, &execErr) {
			os.Exit(execErr.ExitCode())
		}
		os.Exit(1)
	}
}

const goFlag = "go"
const goimportsFlag = "goimports"
const toolsfileFlag = "tools_file"
const toolsdirFlag = "tools_directory"
const configfileFlag = "config_file"
const basedirFlag = "base_dir"
const buildFlagsFlag = "build_flags"
const verboseFlag = "verbose"

func init() {
	rootCmd.PersistentFlags().String(goFlag, "", "The \"go\" executable to use.")
	viper.BindPFlag(goFlag, rootCmd.PersistentFlags().Lookup(goFlag))

	rootCmd.PersistentFlags().String(goimportsFlag, "", "the \"goimports\" executable to use.")
	viper.BindPFlag(goimportsFlag, rootCmd.PersistentFlags().Lookup(goimportsFlag))

	rootCmd.PersistentFlags().String(basedirFlag, "", "the base directory for automatically calculating where to put the tools file and tools directory.  Defaults to the directory where your go.mod file is.")
	viper.BindPFlag(basedirFlag, rootCmd.PersistentFlags().Lookup(basedirFlag))

	rootCmd.PersistentFlags().String(toolsfileFlag, "", "the file in which to store tool data. This should end in a \".go\" extenstion so that go's module system picks it up. Defaults to \"tools.go\" in the the base directory.")
	viper.BindPFlag(toolsfileFlag, rootCmd.PersistentFlags().Lookup(toolsfileFlag))

	rootCmd.PersistentFlags().String(toolsdirFlag, "", "the directory where tool binaries are stored.  Defaults to \"_tools\" in the base directory")
	viper.BindPFlag(toolsdirFlag, rootCmd.PersistentFlags().Lookup(toolsdirFlag))

	rootCmd.PersistentFlags().String(buildFlagsFlag, "", "Any build flags to use when adding a new tool. These are stored and used when syncing the tool in the future.")
	viper.BindPFlag(buildFlagsFlag, rootCmd.PersistentFlags().Lookup(buildFlagsFlag))

	rootCmd.PersistentFlags().BoolP(verboseFlag, "v", false, "Increased outout")
	viper.BindPFlag(verboseFlag, rootCmd.PersistentFlags().Lookup(verboseFlag))

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

func makeOptions() []toolbox.Option {
	options := []toolbox.Option{}
	if goOption := viper.GetString(goFlag); goOption != "" {
		options = append(options, toolbox.GoOption(goOption))
	}
	if goimportsOption := viper.GetString(goimportsFlag); goimportsOption != "" {
		options = append(options, toolbox.GoimportsOption(goimportsOption))
	}
	if basedirOption := viper.GetString(basedirFlag); basedirOption != "" {
		options = append(options, toolbox.BasedirOption(basedirOption))
	}
	if toolsfileOption := viper.GetString(toolsfileFlag); toolsfileOption != "" {
		options = append(options, toolbox.ToolsfileOption(toolsfileOption))
	}
	if toolsdirOption := viper.GetString(toolsdirFlag); toolsdirOption != "" {
		options = append(options, toolbox.ToolsdirOption(toolsdirOption))
	}
	if buildFlagsOption := viper.GetString(buildFlagsFlag); buildFlagsOption != "" {
		options = append(options, toolbox.BuildFlagsOption(buildFlagsOption))
	}
	if verboseOption := viper.GetBool(verboseFlag); verboseOption {
		options = append(options, toolbox.LoggerOption(log.New(os.Stdout, "", 0)))
	}

	return options
}
