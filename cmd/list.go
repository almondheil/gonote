/*
Copyright © 2024 almond Heil <contact@almendra.dev>
*/
package cmd

import (
	"fmt"
	"os"
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
		list_notes(list_tags, list_long)
	},
}

// TODO: I need a function that is basically an iterator over all the notes, that'd be cool
// TODO: I could also use a function that reads the frontmatter of a specific note file, as a common one

func list_notes(required_tags []string, long bool) {
	// TODO: GOD DAMN IT, WE NEED CONFIG

	homedir := os.Getenv("HOME")
	notedir := filepath.Join(homedir, "Notes")
	notes, err := common.ListNotes(notedir)
	if err != nil {
		// TODO: remove panics
		panic(err)
	}

	// print them notes!
	for _, note := range notes {
		note_path := filepath.Join(notedir, note)
		matter, err := common.ReadHeader(note_path)
		if err != nil {
			panic(err)
		}

		// Skip this iteration if we don't have the required tags
		if !common.TagsMatch(required_tags, matter.Tags) {
			continue
		}

		// print in long or short form
		if !long {
			fmt.Println(note, "-", matter.Tags)
		} else {
			if matter.Title != "" {
				fmt.Println(matter.Title)
			} else {
				fmt.Println(note)
			}
			fmt.Println("  date:", matter.Date)
			fmt.Println("  tags:", matter.Tags)
		}
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringArrayVarP(&list_tags, "tag", "t", make([]string, 0), "Only show notes with tag(s)")
	listCmd.Flags().BoolVarP(&list_long, "long", "l", false, "Long note listings")
}
