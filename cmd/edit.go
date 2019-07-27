package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/pkg/git"
	"gitlab.com/JanMa/go-pass/pkg/store"
	"gitlab.com/JanMa/go-pass/pkg/util"
)

func editPassword(cmd *cobra.Command, args []string) {
	entry, err := PasswordStore.FindEntry(args[0])
	exitOnError(err)
	exitOnError(entry.Decrypt())
	tmp := util.TmpFile()
	val, err := entry.Value()
	exitOnError(err)
	exitOnError(ioutil.WriteFile(tmp, []byte(val), 0600))
	editor := exec.Command(getEditor(), tmp)
	editor.Stdin = os.Stdin
	editor.Stdout = os.Stdout
	editor.Stderr = os.Stderr
	exitOnError(editor.Start())
	exitOnError(editor.Wait())
	newVal, err := ioutil.ReadFile(tmp)
	exitOnError(err)
	exitOnError(os.Remove(tmp))
	if val == string(newVal) {
		fmt.Println("Password unchanged")
		return
	}
	entry.Insert(string(newVal))
	iD := PasswordStore.FindGpgID(entry.Path)
	recv, err := store.ParseGpgID(iD)
	exitOnError(entry.Encrypt(recv))
	git.AddFile(entry.Path, fmt.Sprintf("Edit %s with %s", args[0], getEditor()))
}

func getEditor() string {
	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		editor = "vi"
	}
	return editor
}
