package cmd

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/pkg/git"
	"gitlab.com/JanMa/go-pass/pkg/store"
	"gitlab.com/JanMa/go-pass/pkg/store/entry"
	"gitlab.com/JanMa/go-pass/pkg/util"
)

func insertPassword(cmd *cobra.Command, args []string) {
	e, err := PasswordStore.FindEntry(args[0])
	if err == nil && !ForceInsert &&
		!util.YesNo(fmt.Sprintf("An entry already exists for %s. Overwrite it?", args[0])) {
		os.Exit(1)
	}
	if e == nil {
		e = entry.New(args[0], PasswordStore.Path+"/"+args[0]+".gpg")
		PasswordStore.InsertEntry(e)
	}
	var pass string
	if MultiLine {
		fmt.Printf("Enter contents of %s and press Ctrl+D when finished:\n\n", args[0])
		pass, err = readMultiLine()
	} else {
		pass, err = enterPassword(args[0])
	}
	exitOnError(err)
	e.Insert(pass)
	iD := PasswordStore.FindGpgID(e.Path)
	recv, _ := store.ParseGpgID(iD)
	exitOnError(e.Encrypt(recv))
	git.AddFile(e.Path, fmt.Sprintf("Add given password for %s to store.", args[0]))
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

func enterPassword(name string) (string, error) {
	fmt.Printf("Enter password for %s: ", name)
	pass, e := readPassword()
	if e != nil {
		return "", e
	}
	fmt.Printf("Retype password for %s: ", name)
	passAgain, e := readPassword()
	if e != nil {
		return "", e
	}
	if pass != passAgain {
		return "", fmt.Errorf("the entered passwords do not match")
	}
	return pass + "\n", nil
}

func readMultiLine() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	pass := []string{}
	for scanner.Scan() {
		pass = append(pass, scanner.Text())
	}
	return strings.Join(pass, "\n") + "\n", scanner.Err()
}
