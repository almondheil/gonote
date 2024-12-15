/*
Copyright © 2024 almond Heil <contact@almendra.dev>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/almondheil/gonote/common"
	"github.com/koki-develop/go-fzf"
	"github.com/spf13/cobra"
)

var edit_tags []string

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:                   "edit [-t tag [-t tag ...]]",
	Short:                 "Search for and open a note.",
	Aliases:               []string{"e"},
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: use config to get user note directory
		homedir := os.Getenv("HOME")
		notedir := filepath.Join(homedir, "Notes")
		notes, err := common.FindNotesFiltered(notedir, edit_tags)
		if err != nil {
			panic(err)
		}

		choice_filenames, err := fzf_choose_notes(notes)
		if err != nil {
			panic(err)
		}

		common.EditNotes(notedir, choice_filenames)
	},
}

func fzf_choose_notes(notes []common.Note) ([]string, error) {
	// Put together the filenames of all our notes
	filenames := make([]string, len(notes))
	filenames_tagged := make([]string, len(notes))
	for i, note := range notes {
		// Store normal filename
		filenames[i] = note.Filename

		// Store filename with tags listed
		tag_formatted := fmt.Sprintf("%s %v", note.Filename, note.Matter.Tags)
		filenames_tagged[i] = tag_formatted
	}

	// Create a fuzzy finder and run it on items
	// from https://github.com/koki-develop/go-fzf/blob/main/README.md#using-as-a-library
	f, err := fzf.New()
	if err != nil {
		return nil, err
	}
	idxs, err := f.Find(filenames_tagged, func(i int) string { return filenames_tagged[i] })
	if err != nil {
		return nil, err
	}

	// Put the filenames the user chose into a slice -- do not include tags!
	choices := make([]string, len(idxs))
	for choice_idx, filename_idx := range idxs {
		choices[choice_idx] = filenames[filename_idx]
	}

	return choices, nil
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().StringArrayVarP(&edit_tags, "tags", "t", make([]string, 0), "Only show notes with tag(s)")
}
