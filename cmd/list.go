/*
Copyright © 2024 almond Heil <contact@almendra.dev>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var list_tags []string;
var list_long bool;

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use: "list [-l] [-t tag [-t tag ...]]",
	Short: "List notes and filter by tags.",
	Aliases: []string{"l", "ls"},
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
		fmt.Println("list_tags is", list_tags)
		fmt.Println("list_long is", list_long)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringArrayVarP(&list_tags, "tag", "t", make([]string, 0), "Only show notes with tag(s)")
	listCmd.Flags().BoolVarP(&list_long, "long", "l", false, "Long note listings")
}
