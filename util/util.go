package util

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh/terminal"
)

// PrintLine prints a line of output
func PrintLine(s string) {
	if i := strings.LastIndex(s, ".gpg"); i != -1 {
		fmt.Println(s[:i])
	} else {
		fmt.Println(s)
	}
}

// RundCommand runs a given command and returns the output
func RunCommand(name string, args ...string) []string {
	cmd, err := exec.Command(name, args...).Output()
	if err != nil {
		fmt.Println(err)
		defer os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(cmd))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func getHomeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		defer os.Exit(1)
	}
	return home
}

// GetPasswordStore returns the path to the password store
func GetPasswordStore() string {
	env := os.Getenv("PASSWORD_STORE_DIR")
	if len(env) == 0 {
		env = getHomeDir() + "/.password-store"
	}
	return env
}

// YesNo simple Yes or No dialogue
func YesNo() bool {
	t := terminal.NewTerminal(os.Stdin, "")
	oldState, _ := terminal.MakeRaw(0)
	defer terminal.Restore(0, oldState)
	fmt.Printf(" [y/N] ")
	i, e := t.ReadLine()
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	i = strings.ToLower(i)
	if i == "y" || i == "yes" {
		return true
	}
	return false
}
