package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"gitlab.com/JanMa/go-pass/pkg/store/entry"
)

// Store represents the password store
type Store struct {
	Path       string
	entries    map[string]*entry.Entry
	gpgID      string
	recipients []string
}

// New returns a new Store
func New(path string) (*Store, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("Path %s does not exist\n", path)
	}
	return &Store{
		Path: path,
	}, nil
}

// Fill reads all .gpg files inside the password store from disk
// and parses them into entry.Entry objects
func (s *Store) Fill() error {
	s.entries = make(map[string]*entry.Entry)
	return filepath.Walk(s.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if filepath.Ext(path) == ".gpg" {
			n := strings.TrimSuffix(strings.TrimPrefix(path, s.Path+"/"), ".gpg")
			e := entry.New(n, path)
			s.entries[n] = e
		}
		return nil
	})
}

// GetPasswordStore returns the path to the password store on disk
func GetPasswordStore() (string, error) {
	env := os.Getenv("PASSWORD_STORE_DIR")
	if len(env) == 0 {
		home, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		env = home + "/.password-store"
	}
	f, e := os.Stat(env)
	if os.IsNotExist(e) {
		return "", e
	}
	if !f.IsDir() {
		return "", fmt.Errorf("%s is not a directory", env)
	}
	return env, nil
}

// FindEntry searches for an entry inside the Store and returns it
func (s *Store) FindEntry(e string) (*entry.Entry, error) {
	var result *entry.Entry
	var err error
	result = s.entries[e]
	if result == nil {
		err = fmt.Errorf("%s is not in the Store", e)
	}
	return result, err
}

// InsertEntry adds a new entry.Entry to the Store
func (s *Store) InsertEntry(e *entry.Entry) {
	s.entries[e.Name] = e
}

// DeleteEntry deletes an entry.Entry from the Store
func (s *Store) DeleteEntry(e *entry.Entry) error {
	var err error
	delete(s.entries, e.Name)
	if s.entries[e.Name] != nil {
		err = fmt.Errorf("Could not delete entry")
	}
	return err
}
