/*
Copyright © 2024 almond Heil <contact@almendra.dev>
*/
package cmd

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Notedir     string `yaml:"note_dir"`
	Extension   string `yaml:"extension"`
	Editor      string `yaml:"editor"`
	ListThreads int    `yaml:"list_threads"`
}

var user_cfg Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "gonote",
	Short:        "Simple note tagging and searching tool.",
	SilenceUsage: true,
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

func find_user_config() error {
	// Get our user home directory
	usr, err := user.Current()
	if err != nil {
		return err
	}

	// Valid places to store config
	search_locations := []string{
		".config/gonote/config.yaml",
		".config/gonote/config.yml",
		".gonote.yaml",
		".gonote.yml",
	}

	// Check in each location
	for _, loc := range search_locations {
		conf_path := filepath.Join(usr.HomeDir, loc)
		if Exists(conf_path) {
			// Read the first config we come across, then be done
			err := read_user_config(conf_path)
			if err != nil {
				return err
			}
			return nil
		}
	}

	// if no files exist, raise that as an error
	return os.ErrNotExist
}

// Read user config into the config struct, returning an error if anything goes wrong.
func read_user_config(path string) error {
	dat, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Read yaml into our global struct
	err = yaml.Unmarshal(dat, &user_cfg)
	if err != nil {
		return err
	}

	// Make sure config's notes dir isn't undefined
	if user_cfg.Notedir == "" {
		return errors.New("config: note_dir is not defined")
	} else {
		// interpret any env vars in notedir (such as $HOME, importantly)
		user_cfg.Notedir = os.ExpandEnv(user_cfg.Notedir)
	}

	// If other config options are undefined, set defaults
	if user_cfg.Extension == "" {
		user_cfg.Extension = ".md"
	}
	if user_cfg.Editor == "" {
		user_cfg.Editor = "vim"
	}
	if user_cfg.ListThreads == 0 {
		user_cfg.ListThreads = 4
	}
	return nil
}
