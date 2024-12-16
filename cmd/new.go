/*
Copyright © 2024 almond Heil <contact@almendra.dev>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var new_tags []string
var new_noprompt bool

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:                   "new [-l] [-t tag [-t tag ...]] title",
	Aliases:               []string{"n"},
	Short:                 "Create a new note.",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := find_user_config()
		if err != nil {
			return err
		}

		err = create_templated_note(user_cfg.Notedir, args[0])
		if err != nil {
			return err
		}
		return nil
	},
}

func create_templated_note(notedir string, title string) error {
	// Get the current date and put it into the right formats
	now := time.Now()
	date_long := now.Format("2006-01-02")
	date_compact := now.Format("20060102")

	// Check whether a note with this title exists already
	note_title := date_compact + "-" + title + user_cfg.Extension
	note_path := filepath.Join(notedir, note_title)
	if Exists(note_path) {
		fmt.Println("Note", note_title, "already exists, did you mean `gonote edit`?")
		os.Exit(1)
	}

	// If enabled for whether the user wants to create the note
	if !new_noprompt {
		create, err := confirm_create(note_title)
		if err != nil {
			return err
		} else if !create {
			return nil
		}
	}

	// Create the for the note
	f, err := os.Create(note_path)
	if err != nil {
		return err
	}
	defer f.Close()

	format := "---\ntitle: %s\ndate: %s\ntags: %v\n---\n"
	matter := fmt.Sprintf(format, title, date_long, new_tags)
	_, err = f.WriteString(matter)
	if err != nil {
		return err
	}

	// Open that file in the user's editor
	err = EditNotes(notedir, []string{note_title})
	if err != nil {
		return err
	}

	return nil
}

func confirm_create(title string) (bool, error) {
	scanner := bufio.NewScanner(os.Stdin)

	// Repeatedly print prompt and scan response
	fmt.Printf("Create note %s? [Y/n] ", title)
	for scanner.Scan() {
		// get the response text in lowercase
		response := strings.ToLower(scanner.Text())

		// empty line, line starting with y -> true
		// line starting with n -> false
		// anything else -> repeat loop
		if response == "" || response[0] == 'y' {
			return true, nil
		} else if response[0] == 'n' {
			return false, nil
		}

		// this is actually the print for the next iter bcs the loop is screwy
		fmt.Printf("Create note %s? [Y/n] ", title)
	}

	// TODO: this is the case where scanner stops for some reason, what do I do here?
	return false, nil
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringArrayVarP(&new_tags, "tags", "t", make([]string, 0), "Create note with tag(s)")
	newCmd.Flags().BoolVarP(&new_noprompt, "yes", "y", false, "Create note without prompting")
}
