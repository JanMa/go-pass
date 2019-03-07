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
	Use:   "edit",
	Short: "Insert a new password or edit an existing password using /usr/bin/nano.",
	Args:  cobra.ExactArgs(1),
	Run:   editPassword,
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func editPassword(cmd *cobra.Command, args []string) {
	root := util.GetPasswordStore() + "/" + args[0] + ".gpg"
	tmpfile := os.TempDir() + "/go-pass-tmp"
	if f, e := os.Stat(root); !os.IsNotExist(e) && !f.IsDir() {
		decrypt := exec.Command("gpg", "--quiet", "-o", tmpfile, "-d", root)
		if err := decrypt.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	cmdArgs := "cat \"" + tmpfile + "\" | gpg -e " + getRecepientOpts() + " -o " + strings.ReplaceAll(root, " ", `\ `) + " --quiet --yes --compress-algo=none --no-encrypt-to"
	gpg := exec.Command("bash", "-c", cmdArgs)
	edit := os.Getenv("EDITOR")
	if len(edit) == 0 {
		edit = "nano"
	}
	editor := exec.Command(edit, tmpfile)
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

	os.MkdirAll(filepath.Dir(root), 0700)
	if err := gpg.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := os.Remove(tmpfile); err != nil {
		fmt.Println(err)
	}

}
