package cmd

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/pkg/git"
	"gitlab.com/JanMa/go-pass/util"
)

func insertPassword(cmd *cobra.Command, args []string) {
	root := util.GetPasswordStore() + "/" + args[0] + ".gpg"
	if f, e := os.Stat(root); !os.IsNotExist(e) && !f.IsDir() {
		if ForceInsert || util.YesNo(fmt.Sprintf("An entry already exists for %s. Overwrite it?", args[0])) {
			exitOnError(os.Remove(root))
		} else {
			os.Exit(1)
		}
	}
	if MultiLine {
		fmt.Printf("Enter contents of %s and press Ctrl+D when finished:\n\n", args[0])
		encryptMultiLine(root)
	} else {
		encryptPassword(enterPassword(args[0])+"\n", root)
	}
	git.AddFile(root, fmt.Sprintf("Add given password for %s to store.", args[0]))
}

func getRecepientOptsArray() []string {
	opts := []string{}
	idFile := util.GetPasswordStore() + "/.gpg-id"
	i, e := os.Open(idFile)
	exitOnError(e)
	s := bufio.NewScanner(i)
	for s.Scan() {
		opts = append(opts, "-r")
		opts = append(opts, s.Text())
	}
	return opts
}

func readPassword() (string, error) {
	t := terminal.NewTerminal(os.Stdin, "")
	oldState, _ := terminal.MakeRaw(0)
	defer terminal.Restore(0, oldState)
	if Echo {
		p, e := t.ReadLine()
		return p, e
	}
	p, e := t.ReadPassword("")
	return p, e
}

func enterPassword(name string) string {
	fmt.Printf("Enter password for %s: ", name)
	pass, e := readPassword()
	exitOnError(e)
	fmt.Printf("Retype password for %s: ", name)
	passAgain, e := readPassword()
	exitOnError(e)
	if pass != passAgain {
		fmt.Println("Error: the entered passwords do not match.")
		os.Exit(1)
	}
	return pass
}

func encryptPassword(pass, file string) {
	gpg := exec.Command("gpg",
		"-e", "-o", strings.ReplaceAll(file, " ", `\ `),
		"--quiet", "--yes", "--compress-algo=none", "--no-encrypt-to")
	for _, r := range getRecepientOptsArray() {
		gpg.Args = append(gpg.Args, r)
	}
	stdin, err := gpg.StdinPipe()
	exitOnError(err)
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, pass)
	}()
	os.MkdirAll(filepath.Dir(file), 0755)
	exitOnError(gpg.Run())
}

func encryptMultiLine(file string) {
	scanner := bufio.NewScanner(os.Stdin)
	pass := []string{}
	for scanner.Scan() {
		pass = append(pass, scanner.Text())
	}
	exitOnError(scanner.Err())
	encryptPassword(strings.Join(pass, "\n")+"\n", file)
}
