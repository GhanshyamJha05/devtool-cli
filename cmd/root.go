package cmd

import (
	"fmt"
	"os"

	"github.com/GhanshyamJha05/devtool-cli/internal/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Verbose is a global flag accessible by all subcommands for debug logging.
var Verbose bool

var rootCmd = &cobra.Command{
	Use:   "devtool",
	Short: "A production-quality CLI developer tool",
	Long: `devtool-cli is a modular command-line application that automates
common developer tasks like fetching APIs, formatting JSON, and organizing files.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".devtool-cli")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if Verbose {
			utils.PrintDebug(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
		}
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "enable verbose/debug output")
}
