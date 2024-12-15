/*
Copyright © 2024 almond Heil <contact@almendra.dev>
*/

package common

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/adrg/frontmatter"
)

type NoteFrontmatter struct {
	Title string   `yaml:"title"`
	Date  string   `yaml:"date"`
	Tags  []string `yaml:"tags"`
}

type Note struct {
	Filename    string
	Frontmatter NoteFrontmatter
}

func ListNotes(path string) ([]string, error) {
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

func ReadHeader(path string) (NoteFrontmatter, error) {
	var matter NoteFrontmatter

	// Attempt to read the file contents
	dat, err := os.ReadFile(path)
	if err != nil {
		return matter, err
	}

	// Try to parse the frontmatter of the file
	contents := string(dat)
	_, err = frontmatter.Parse(strings.NewReader(contents), &matter)
	if err != nil {
		return matter, err
	}

	return matter, nil
}

func TagsMatch(required []string, check []string) bool {
	// just loop over the required tags and be slow, because why not
	for _, tag := range required {
		if !slices.Contains(check, tag) {
			return false
		}
	}

	return true
}

func FindNotesFiltered(notedir string, required_tags []string) ([]Note, error) {
	found_notes := make([]Note, 0)

	notes, err := ListNotes(notedir)
	if err != nil {
		return nil, err
	}

	// collect the notes into a slice
	for _, note_filename := range notes {
		note_path := filepath.Join(notedir, note_filename)
		matter, err := ReadHeader(note_path)
		if err != nil {
			return nil, err
		}

		// Skip this iteration if we don't have the required tags
		if !TagsMatch(required_tags, matter.Tags) {
			continue
		}

		found_notes = append(found_notes, Note{Filename: note_filename, Frontmatter: matter})
	}

	return found_notes, nil
}
