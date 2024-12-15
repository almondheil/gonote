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
		list_note_tags()
	},
}

func list_note_tags() {
	homedir := os.Getenv("HOME")
	notedir := filepath.Join(homedir, "Notes")
	notes, err := common.ListNotes(notedir)
	if err != nil {
		// TODO: remove panics
		panic(err)
	}

	// store a set of the tags we've seen, and go through all notes to find them
	var found_tags = make(map[string]bool)
	for _, note := range notes {
		note_path := filepath.Join(notedir, note)
		matter, err := common.ReadHeader(note_path)
		if err != nil {
			panic(err)
		}

		for _, t := range matter.Tags {
			found_tags[t] = true
		}
	}

	// Print out the tags we found in alphabetical order
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
