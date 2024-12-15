/*
Copyright © 2024 almond Heil <contact@almendra.dev>
*/
package cmd

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/almondheil/gonote/common"
	"github.com/spf13/cobra"
)

var list_tags []string
var list_long bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:                   "list [-l] [-t tag [-t tag ...]]",
	Short:                 "List notes and filter by tags.",
	Aliases:               []string{"l", "ls"},
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: use config to get user home dir
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}
		notedir := filepath.Join(usr.HomeDir, "Notes")
		notes, err := common.FindNotesFiltered(notedir, list_tags)
		if err != nil {
			panic(err)
		}

		// Print out each note that has the correct tags
		for _, n := range notes {
			print_note_info(n, list_long)
		}
	},
}

func print_note_info(note common.Note, long bool) {
	// print in long or short form
	if !long {
		fmt.Println(note.Filename, "-", note.Matter.Tags)
	} else {
		if note.Matter.Title != "" {
			fmt.Println(note.Matter.Title)
		} else {
			fmt.Println(note.Filename)
		}
		fmt.Println("  date:", note.Matter.Date)
		fmt.Println("  tags:", note.Matter.Tags)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringArrayVarP(&list_tags, "tag", "t", make([]string, 0), "Only show notes with tag(s)")
	listCmd.Flags().BoolVarP(&list_long, "long", "l", false, "Long note listings")
}
