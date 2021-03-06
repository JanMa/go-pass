package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
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
		return nil, fmt.Errorf("path %s does not exist", path)
	}
	rec, _ := ParseGpgID(path + "/.gpg-id")
	return &Store{
		Path:       path,
		gpgID:      path + "/.gpg-id",
		recipients: rec,
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

// FindEntries searches for entries inside the Store and returns
// a map containing the results
func (s *Store) FindEntries(e string) (map[string]*entry.Entry, error) {
	m := make(map[string]*entry.Entry)
	r := regexp.MustCompile("^" + e + "$")
	var err error
	for k := range s.entries {
		if r.MatchString(k) {
			m[k] = s.entries[k]
		}
	}
	if len(m) == 0 {
		err = fmt.Errorf("Found no matching entires for %s", e)
	}
	return m, err
}

// SortEntries returns a sorted array of entry names
// for a given map of entries
func SortEntries(entries map[string]*entry.Entry) []string {
	names := []string{}
	for e := range entries {
		names = append(names, e)
	}
	sort.Strings(names)
	return names
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

// FindGpgID traverses the store to find the next matching .gpg-id
// file for a given path
func (s *Store) FindGpgID(path string) string {
	dirs := strings.Split(strings.Trim(path, "/"), "/")
	root := strings.Split(strings.Trim(s.Path, "/"), "/")
	if l := len(dirs); filepath.Ext(dirs[l-1]) == ".gpg" {
		dirs = dirs[:(l - 1)]
	}
	for i := len(dirs); i >= len(root); i-- {
		p := "/" + strings.Join(dirs[:i], "/") + "/.gpg-id"
		if f, e := os.Stat(p); !os.IsNotExist(e) && !f.IsDir() {
			return p
		}
	}
	return s.Path + "/.gpg-id"
}

// ParseGpgID takes the path to a .gpg-id file and returns
// it's entries
func ParseGpgID(path string) ([]string, error) {
	gpgID, e := ioutil.ReadFile(path)
	if e != nil {
		return nil, e
	}
	return strings.Split(strings.Trim(string(gpgID), "\n"), "\n"), nil
}

// ShowAll returns a map of all entries in the store
func (s *Store) ShowAll() map[string]*entry.Entry {
	return s.entries
}

//GetGpgKeys returns the gpg keys of an array of recipients
func GetGpgKeys(recipients []string) ([]string, error) {
	re := regexp.MustCompile(`sub:[^:]*:[^:]*:[^:]*:([^:]*):[^:]*:[^:]*:[^:]*:[^:]*:[^:]*:[^:]*:[a-zA-Z]*e[a-zA-Z]*:.*`)
	gpg := exec.Command("gpg", "--list-keys", "--with-colons")
	for _, r := range recipients {
		gpg.Args = append(gpg.Args, r)
	}
	k, e := gpg.Output()
	if e != nil {
		return nil, e
	}
	match := re.FindAllSubmatch(k, -1)

	gpgKeys := []string{}
	for _, m := range match {
		gpgKeys = append(gpgKeys, string(m[1]))
	}

	return gpgKeys, nil
}
