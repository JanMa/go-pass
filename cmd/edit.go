// Copyright © 2019 Jan Martens
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
