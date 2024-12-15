/*
Copyright © 2024 almond Heil <contact@almendra.dev>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use: "new [-l] [-t tag [-t tag ...]] [title]",
	Aliases: []string{"n"},
	Short: "Create a new note.",
	DisableFlagsInUseLine: true,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("new called")
		fmt.Println("args are", args)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringArrayP("tags", "t", make([]string, 0), "Create note with tag(s)")
	newCmd.Flags().BoolP("yes", "y", false, "Create note without prompting")
}
