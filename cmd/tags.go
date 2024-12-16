/*
Copyright © 2024 almond Heil <contact@almendra.dev>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

var tags_show_count bool

// tagsCmd represents the tags command
var tagsCmd = &cobra.Command{
	Use:     "tags",
	Short:   "List tags across existing notes.",
	Aliases: []string{"t"},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := find_user_config()
		if err != nil {
			return err
		}

		// Get all notes in that directory (regardless of tag)
		homedir := os.Getenv("HOME")
		notedir := filepath.Join(homedir, "Notes")
		notes, err := FindNotesFiltered(notedir, []string{})
		if err != nil {
			return err
		}

		// Print out the unique tags across all notes
		print_note_tags(notes)
		return nil
	},
}

func print_note_tags(notes []Note) {
	// Store a set of the tags we've seen, and go through all notes to find them
	var found_tags = make(map[string]int)
	for _, note := range notes {
		for _, t := range note.Matter.Tags {
			found_tags[t] += 1
		}
	}

	// Pull key-value pairs into a struct together
	type KVPair struct {
		s string
		i int
	}
	pairs := make([]KVPair, len(found_tags))
	i := 0
	for k := range found_tags {
		pairs[i].s = k
		pairs[i].i = found_tags[k]
		i++
	}

	// Sort and print depending on whether we are showing the count or not
	if tags_show_count {
		// sort by count ascending
		slices.SortFunc(pairs, func(a, b KVPair) int {
			return a.i - b.i
		})

		// print out both numbers and titles
		for _, pair := range pairs {
			fmt.Printf("%d %s\n", pair.i, pair.s)
		}
	} else {
		// sort by alphabetic ascending
		slices.SortFunc(pairs, func(a, b KVPair) int {
			return strings.Compare(a.s, b.s)
		})

		// print out just the titles (strings)
		for _, pair := range pairs {
			fmt.Println(pair.s)
		}
	}
}

func init() {
	rootCmd.AddCommand(tagsCmd)

	tagsCmd.Flags().BoolVar(&tags_show_count, "count", false, "Sort and print by usage count.")
}
