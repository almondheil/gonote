/*
Copyright © 2024 almond Heil <contact@almendra.dev>
*/

package common

import (
	"os"
	"strings"

	"github.com/adrg/frontmatter"
)

type NoteFrontmatter struct {
	Title string   `yaml:"title"`
	Date  string   `yaml:"date"`
	Tags  []string `yaml:"tags"`
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
