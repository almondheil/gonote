/*
Copyright © 2024 almond Heil <contact@almendra.dev>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [-l] [-t tag [-t tag ...]]",
	Short: "List notes and filter by tags.",
	Aliases: []string{"l", "ls"},
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringArrayP("tag", "t", make([]string, 0), "Only show notes with tag(s)")
	listCmd.Flags().BoolP("long", "l", false, "Long note listings")
}
