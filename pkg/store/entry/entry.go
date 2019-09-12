package entry

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	// "strings"
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
	d, err := exec.Command("gpg", "-dq", e.Path).CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(d))
	}
	if len(d) > 0 && err == nil {
		e.value = string(d)
	}
	return err
}

// Encrypt writes encrypted entry to disk
func (e *Entry) Encrypt(recipients []string) error {
	gpg := exec.Command("gpg",
		"-e", "-o", e.Path,
		"--quiet", "--yes", "--compress-algo=none", "--no-encrypt-to")
	for _, r := range recipients {
		gpg.Args = append(gpg.Args, "-r")
		gpg.Args = append(gpg.Args, r)
	}
	// We need to trust all keys when running in test mode.
	// This is rather ugly but there is no way around it.
	if testRunning {
		gpg.Args = append(gpg.Args, "--trust-model", "always")
	}
	stdin, err := gpg.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, e.value)
	}()
	os.MkdirAll(filepath.Dir(e.Path), 0755)
	out, err := gpg.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
	}
	return err
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

// GetKeys returns a string containing all current gpg keys
// used to encrypt the Entry
func (e *Entry) GetKeys() (string, error) {
	if _, e := os.Stat(e.Path); os.IsNotExist(e) {
		return "", fmt.Errorf("Entry not encrypted")
	}
	k, err := exec.Command("gpg", "-v", "-d", "--list-only", "--keyid-format", "long", e.Path).CombinedOutput()
	return string(k), err
}
