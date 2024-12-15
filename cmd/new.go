/*
Copyright © 2024 almond Heil <contact@almendra.dev>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var new_tags []string;
var new_noprompt bool;

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use: "new [-l] [-t tag [-t tag ...]] [title]",
	Aliases: []string{"n"},
	Short: "Create a new note.",
	DisableFlagsInUseLine: true,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("new called")
		fmt.Println("args is", args)
		fmt.Println("new_tags is", new_tags)
		fmt.Println("new_noprompt is", new_noprompt)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringArrayVarP(&new_tags, "tags", "t", make([]string, 0), "Create note with tag(s)")
	newCmd.Flags().BoolVarP(&new_noprompt, "yes", "y", false, "Create note without prompting")
}
