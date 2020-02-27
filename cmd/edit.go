package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/pkg/git"
	"gitlab.com/JanMa/go-pass/pkg/store"
	"gitlab.com/JanMa/go-pass/pkg/store/entry"
	"gitlab.com/JanMa/go-pass/pkg/util"
)

func editPassword(cmd *cobra.Command, args []string) {
	e, exists := PasswordStore.FindEntry(args[0])
	tmp, err := util.TmpFile()
	exitOnError(err)
	if exists == nil {
		exitOnError(e.Decrypt())
		val, err := e.Value()
		exitOnError(err)
		exitOnError(ioutil.WriteFile(tmp, []byte(val), 0600))
	} else {
		fmt.Println("Creating new entry", args[0])
		e = entry.New(args[0], PasswordStore.Path+"/"+args[0]+".gpg")
		PasswordStore.InsertEntry(e)
	}
	editor := exec.Command(getEditor(), tmp)
	editor.Stdin = os.Stdin
	editor.Stdout = os.Stdout
	editor.Stderr = os.Stderr
	exitOnError(editor.Start())
	exitOnError(editor.Wait())
	newVal, err := ioutil.ReadFile(tmp)
	exitOnError(err)
	exitOnError(os.Remove(tmp))
	if val, _ := e.Value(); val == string(newVal) {
		fmt.Println("Password unchanged")
		return
	}
	e.Insert(string(newVal))
	iD := PasswordStore.FindGpgID(e.Path)
	recv, err := store.ParseGpgID(iD)
	exitOnError(e.Encrypt(recv))
	git.AddFile(e.Path, fmt.Sprintf("Edit %s with %s", args[0], getEditor()))
}

func getEditor() string {
	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		editor = "vi"
	}
	return editor
}
