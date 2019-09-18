package entry

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gitlab.com/JanMa/go-pass/pkg/gpg"
)

// Entry represents an entry inside the password store
type Entry struct {
	Name  string
	Path  string
	value string
}

var (
	testRunning = false
)

// Decrypt reads encrypted entry from disk
func (e *Entry) Decrypt() error {
	if _, err := os.Stat(e.Path); os.IsNotExist(err) {
		return fmt.Errorf("%s is not in the password store", e.Name)
	}
	d, err := gpg.Decrypt(e.Path)
	if err != nil {
		return fmt.Errorf(d)
	}
	e.value = d
	return err
}

// Encrypt writes encrypted entry to disk
func (e *Entry) Encrypt(recipients []string) error {
	return gpg.Encrypt(e.Path, e.value, recipients, testRunning)
}

// Show prints the entry value
func (e *Entry) Show() error {
	if len(e.value) > 0 {
		fmt.Print(e.value)
		return nil
	}
	return fmt.Errorf("Entry not decrypted")
}

// New returns a new Entry
func New(name, path string) *Entry {
	return &Entry{
		Name: name,
		Path: path,
	}
}

// Value returns the entry value
func (e *Entry) Value() (string, error) {
	if len(e.value) > 0 {
		return e.value, nil
	}
	return "", fmt.Errorf("Entry not decrypted")
}

// Insert writes new value to the entry
func (e *Entry) Insert(value string) {
	e.value = value
}

// Delete deletes an encrypted entry from disk
func (e *Entry) Delete() error {
	err := os.RemoveAll(e.Path)
	if err != nil {
		return err
	}
	dir := filepath.Dir(e.Path)
	empty, err := isEmpty(dir)
	if err != nil {
		return err
	}
	if empty {
		err = os.RemoveAll(dir)
	}
	return err
}

// https://stackoverflow.com/a/30708914
func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}
