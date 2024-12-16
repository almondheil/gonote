/*
Copyright © 2024 almond Heil <contact@almendra.dev>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

type Config struct {
	Notedir   string `yaml:"note_dir"`
	Extension string `yaml:"extension"`
	Editor    string `yaml:"editor"`
}

var user_cfg Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gonote",
	Short: "Simple note tagging and searching tool.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// nothing here, we don't do anything when gonote is run without a subcommand
}
