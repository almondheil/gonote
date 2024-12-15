/*
Copyright © 2024 almond Heil <contact@almendra.dev>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [-t tag [-t tag ...]]",
	Short: "Search for and open a note.",
	Aliases: []string{"e"},
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("edit called")
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().StringArrayP("tags", "t", make([]string, 0), "Only show notes with tag(s)")
}
