/*
Copyright © 2024 almond Heil <contact@almendra.dev>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/koki-develop/go-fzf"
	"github.com/spf13/cobra"
)

var edit_tags []string

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:                   "edit [-t tags [-t tags ...]]",
	Short:                 "Search for and open a note.",
	Aliases:               []string{"e"},
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := find_user_config()
		if err != nil {
			return err
		}

		notedir := user_cfg.Notedir
		notes, err := FindNotesFiltered(notedir, edit_tags)
		if err != nil {
			return err
		}

		// if 0 or 1 notes are found, special behavior happens.
		// - 0 notes is a exit and note to the user
		// - 1 note means we just open it right away
		if len(notes) == 0 {
			fmt.Fprintln(os.Stderr, "No notes to edit")
			return nil
		} else if len(notes) == 1 {
			choices := make([]string, 1)
			choices[0] = notes[0].Filename
			EditNotes(notedir, choices)
			return nil
		}

		choice_filenames, err := fzf_choose_notes(notes)
		if err != nil {
			return err
		}

		EditNotes(notedir, choice_filenames)
		return nil
	},
}

func fzf_choose_notes(notes []Note) ([]string, error) {
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

	editCmd.Flags().StringSliceVarP(&edit_tags, "tags", "t", make([]string, 0), "Only show notes with tag(s)")
}
