package util

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

// Prints a line of output
func PrintLine(s string) {
	if i := strings.LastIndex(s, ".gpg"); i != -1 {
		fmt.Println(s[:i])
	} else {
		fmt.Println(s)
	}
}

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

func GetHomeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		defer os.Exit(1)
	}
	return home
}
