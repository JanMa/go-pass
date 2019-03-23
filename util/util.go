package util

import (
	"fmt"
	"math/rand"
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

// RunCommand runs a given command and returns the output
func RunCommand(name string, args ...string) []string {
	cmd, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		fmt.Println(string(cmd))
		defer os.Exit(1)
	}
	return strings.Split(string(cmd), "\n")
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

// RandomString returns a random string of [a-zA-Z0-1] of the given length
func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

// TmpFile generates the path to a new temporary file
func TmpFile() string {
	os.Mkdir(os.TempDir()+"/go-pass", 0700)
	return os.TempDir() + "/go-pass/" + RandomString(8)
}
