/*
Copyright © 2024 almond Heil <contact@almendra.dev>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/almondheil/gonote/common"
	"github.com/spf13/cobra"
)

// tagsCmd represents the tags command
var tagsCmd = &cobra.Command{
	Use:     "tags",
	Short:   "List tags across existing notes.",
	Aliases: []string{"t"},
	Run: func(cmd *cobra.Command, args []string) {
		// Get all notes in that directory (regardless of tag)
		homedir := os.Getenv("HOME")
		notedir := filepath.Join(homedir, "Notes")
		notes, err := common.FindNotesFiltered(notedir, []string{})
		if err != nil {
			// TODO: remove panics
			panic(err)
		}

		// Print out the unique tags across all notes
		print_note_tags(notes)

	},
}

func print_note_tags(notes []common.Note) {
	// Store a set of the tags we've seen, and go through all notes to find them
	var found_tags = make(map[string]bool)
	for _, note := range notes {
		for _, t := range note.Frontmatter.Tags {
			found_tags[t] = true
		}
	}

	// Get a slice with every tag we've seen
	// from https://stackoverflow.com/a/27848197
	keys := make([]string, len(found_tags))
	i := 0
	for k := range found_tags {
		keys[i] = k
		i++
	}

	// Sort alphabetically before printing
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println(k)
	}
}

func init() {
	rootCmd.AddCommand(tagsCmd)
}
