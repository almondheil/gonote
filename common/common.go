/*
Copyright © 2024 almond Heil <contact@almendra.dev>
*/

package common

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"syscall"

	"github.com/adrg/frontmatter"
)

type note_frontmatter struct {
	Title string   `yaml:"title"`
	Date  string   `yaml:"date"`
	Tags  []string `yaml:"tags"`
}

type Note struct {
	Filename string
	Matter   note_frontmatter
}

func list_notes(path string) ([]string, error) {
	var files []string

	// Open directory
	dir_entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Append all the entries to the array
	for _, entry := range dir_entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

func tags_match(required []string, check []string) bool {
	// just loop over the required tags and be slow, because why not
	for _, tag := range required {
		if !slices.Contains(check, tag) {
			return false
		}
	}

	return true
}

func EditNotes(notedir string, note_titles []string) error {
	// Locate editor binary
	// TODO: open the note with the configured editor (FOR NOW, JUST VIM)
	// TODO: what if the editor is more than one word? then it's more than one arg and hoooo boy
	binary, err := exec.LookPath("vim")
	if err != nil {
		return err
	}

	// Create the args needed to run vim
	passed_args := make([]string, 1+len(note_titles))
	passed_args[0] = "vim"
	for i, name := range note_titles {
		passed_args[i+1] = filepath.Join(notedir, name)
	}

	// Exec the new process we need, replacing the go process
	env := os.Environ()
	err = syscall.Exec(binary, passed_args, env)
	return err // if we hit this, exec failed anyway
}

func FindNotesFiltered(notedir string, required_tags []string) ([]Note, error) {
	found_notes := make([]Note, 0)

	notes, err := list_notes(notedir)
	if err != nil {
		return nil, err
	}

	// collect the notes into a slice
	for _, note_filename := range notes {
		note_path := filepath.Join(notedir, note_filename)
		matter, err := ReadMatter(note_path)
		if err != nil {
			return nil, err
		}

		// Skip this iteration if we don't have the required tags
		if !tags_match(required_tags, matter.Tags) {
			continue
		}

		found_notes = append(found_notes, Note{Filename: note_filename, Matter: matter})
	}

	return found_notes, nil
}

func ReadMatter(path string) (note_frontmatter, error) {
	var matter note_frontmatter

	// Attempt to read the file contents
	// TODO: Can I avoid reading the entire file? I need at most 5 lines of it,
	// or it would be even better if I could read until the end of the frontmatter
	dat, err := os.ReadFile(path)
	if err != nil {
		return matter, err
	}

	// Try to parse the frontmatter of the file
	// TODO: I bet converting to string is expensive too, I don't like this because it gets called per-note per-command
	contents := string(dat)
	_, err = frontmatter.Parse(strings.NewReader(contents), &matter)
	if err != nil {
		return matter, err
	}

	return matter, nil
}
