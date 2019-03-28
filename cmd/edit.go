package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/util"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:                   "edit pass-name",
	Short:                 "Insert a new password or edit an existing password using " + getEditor() + ".",
	Args:                  cobra.ExactArgs(1),
	Run:                   editPassword,
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func editPassword(cmd *cobra.Command, args []string) {
	root := util.GetPasswordStore() + "/" + args[0] + ".gpg"
	tmpfile := util.TmpFile()
	if f, e := os.Stat(root); !os.IsNotExist(e) && !f.IsDir() {
		decrypt := exec.Command("gpg", "--quiet", "-o", tmpfile, "-d", root)
		if err := decrypt.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	gpg := exec.Command("gpg2",
		"-e", "-o", strings.ReplaceAll(root, " ", `\ `),
		"--quiet", "--yes", "--compress-algo=none", "--no-encrypt-to")
	for _, r := range getRecepientOptsArray() {
		gpg.Args = append(gpg.Args, r)
	}
	gpg.Args = append(gpg.Args, tmpfile)
	editor := exec.Command(getEditor(), tmpfile)
	editor.Stdin = os.Stdin
	editor.Stdout = os.Stdout
	editor.Stderr = os.Stderr
	if err := editor.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := editor.Wait(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.MkdirAll(filepath.Dir(root), 0755)
	if err := gpg.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := os.Remove(tmpfile); err != nil {
		fmt.Println(err)
	}

}

func getEditor() string {
	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		editor = "vi"
	}
	return editor
}
