package gpg

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	// "strings"
)

// Decrypt gpg encrypted file from disk
func Decrypt(path string) (string, error) {
	d, err := exec.Command("gpg", "-dq", path).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(string(d))
	}
	return string(d), err
}

// Encrypt writes encrypted value to disk
func Encrypt(path, value string, recipients []string, testRunning bool) error {
	gpg := exec.Command("gpg",
		"-e", "-o", path,
		"--quiet", "--yes", "--compress-algo=none", "--no-encrypt-to")

	// If PASSWORD_STORE_ARMOR is set, enable ASCII encoded output
	if len(os.Getenv("PASSWORD_STORE_ARMOR")) > 0 {
		gpg.Args = append(gpg.Args, "--armor")
	}

	for _, r := range recipients {
		gpg.Args = append(gpg.Args, "-r", r)
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
		io.WriteString(stdin, value)
	}()
	os.MkdirAll(filepath.Dir(path), 0755)
	out, err := gpg.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
	}
	return err
}

// GetKeys returns a string containing all current gpg keys
// used to encrypt a file
func GetKeys(path string) (string, error) {
	if _, e := os.Stat(path); os.IsNotExist(e) {
		return "", fmt.Errorf("Entry not encrypted")
	}
	k, err := exec.Command("gpg", "-v", "-d", "--list-only", "--keyid-format", "long", path).CombinedOutput()
	return string(k), err
}
