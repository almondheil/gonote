/*
Copyright © 2024 almond Heil <contact@almendra.dev>
*/

package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"slices"
	"sort"
	"syscall"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v3"
)

type Frontmatter struct {
	Title string   `yaml:"title"`
	Date  string   `yaml:"date"`
	Tags  []string `yaml:"tags,flow"`
}

type Note struct {
	Filename string
	Matter   Frontmatter
}

func list_notes(path string) ([]string, error) {
	// Open directory
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Use Readdirnames to JUST get the names of all the files.
	// This is really fast, but doesn't let us check if they're ok or not.
	note_names, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	return note_names, nil
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
	// TODO: what if the editor is more than one word? then it's more than one arg and hoooo boy
	binary, err := exec.LookPath(user_cfg.Editor)
	if err != nil {
		return err
	}

	// Create the args needed to run vim
	passed_args := make([]string, 1+len(note_titles))
	passed_args[0] = user_cfg.Editor
	for i, name := range note_titles {
		passed_args[i+1] = filepath.Join(notedir, name)
	}

	// Exec the new process we need, replacing the go process
	env := os.Environ()
	err = syscall.Exec(binary, passed_args, env)
	return err // if we hit this, exec failed anyway
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

type ReadNoteResult struct {
	val Note
	err error
}

func read_note_worker(notedir string, jobs <-chan string, results chan<- ReadNoteResult) {
	for filename := range jobs {
		// Create an empty result that we'll build up
		var res ReadNoteResult
		res.val.Filename = filename

		// Open the file for reading
		path := filepath.Join(notedir, filename)
		file, err := os.Open(path)
		if err != nil {
			res.err = err
			results <- res
			continue
		}

		// Attempt to read its YAML frontmatter
		_, err = frontmatter.Parse(file, &res.val.Matter)
		if err != nil {
			error_str := fmt.Sprintf("reading frontmatter of %s: %s", path, err.Error())
			res.err = errors.New(error_str)
			results <- res
			continue
		}

		// Close the file before we do another iteration
		file.Close()

		// Put our non-error note into the results channel
		res.err = nil
		results <- res
	}
}

func FindNotesFiltered(notedir string, required_tags []string) ([]Note, error) {
	// Get a list of all the notes
	notes, err := list_notes(notedir)
	if err != nil {
		return nil, err
	}

	// We're gonna read the notes in parallel, so create channels to communicate
	num_notes := len(notes)
	jobs := make(chan string, num_notes)
	results := make(chan ReadNoteResult, num_notes)

	// TODO: 4 is an arbitrary number of threads to create. idk man
	for w := 0; w < 4; w++ {
		go read_note_worker(notedir, jobs, results)
	}

	// Send in the needed filenames, letting the workers start
	for _, filename := range notes {
		jobs <- filename
	}
	close(jobs)

	// Collect the results from each goroutine
	found_notes := make([]Note, num_notes)
	failed := false
	for i := 0; i < num_notes; i++ {
		res := <-results

		// Print an error if the goroutine failed
		if res.err != nil {
			fmt.Fprintf(os.Stderr, "note not read: %v\n", err)
			failed = true
		}

		// Store the note in found_notes if the tags match
		if tags_match(required_tags, res.val.Matter.Tags) {
			found_notes = append(found_notes, res.val)
		}
	}

	// If any goroutines failed, return that as an error after collecting all results
	if failed {
		return nil, errors.New("not all notes could be read")
	}

	// Sort our list of found notes by filename, then return them
	sort.Slice(found_notes, func(i, j int) bool {
		return found_notes[i].Filename < found_notes[j].Filename
	})
	return found_notes, nil
}
