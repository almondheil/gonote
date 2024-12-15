/*
Copyright © 2024 almond Heil <contact@almendra.dev>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// tagsCmd represents the tags command
var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List tags across existing notes.",
	Aliases: []string{"t"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tags called")
	},
}

func init() {
	rootCmd.AddCommand(tagsCmd)
}
