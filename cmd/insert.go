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
	"bufio"
	"fmt"
	"gitlab.com/JanMa/go-pass/util"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// insertCmd represents the insert command
var (
	insertCmd = &cobra.Command{
		Use:   "insert",
		Args:  cobra.ExactArgs(1),
		Short: "Insert new password.",
		Long: `Insert new password. Optionally, echo the password back to the console
during entry. Or, optionally, the entry may be multiline. Prompt before
overwriting existing password unless forced.`,
		Run: insertPassword,
	}
	Echo      bool
	MultiLine bool
	Force     bool
)

func init() {
	rootCmd.AddCommand(insertCmd)

	insertCmd.Flags().BoolVarP(&Echo, "echo", "e", false, "Echo password back to console")
	insertCmd.Flags().BoolVarP(&MultiLine, "multiline", "m", false, "Multiline input")
	insertCmd.Flags().BoolVarP(&Force, "force", "f", false, "Overwrite existing password without prompt")
}

func insertPassword(cmd *cobra.Command, args []string) {
	root := util.GetPasswordStore() + "/" + args[0] + ".gpg"
	if f, e := os.Stat(root); !os.IsNotExist(e) && !f.IsDir() {
		fmt.Printf("An entry already exists for %s. Overwrite it?", args[0])
		if util.YesNo() {
			if err := os.Remove(root); err != nil {
				fmt.Println(err)
			}
		} else {
			os.Exit(1)
		}
	}
	if MultiLine {
		fmt.Printf("Enter contents of %s and press Ctrl+D when finished:\n\n", args[0])
		encryptMultiLine(root)
	} else {
		password := enterPassword(args[0])
		encryptPassword(password, root)
	}
}

func getRecepientOpts() string {
	opts := ""
	idFile := util.GetPasswordStore() + "/.gpg-id"
	i, e := os.Open(idFile)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	s := bufio.NewScanner(i)
	for s.Scan() {
		opts += " -r " + s.Text()
	}
	return opts
}

func readPassword() (string, error) {
	if Echo {
		r := bufio.NewReader(os.Stdin)
		p, e := r.ReadString('\n')
		return strings.Trim(p, "\n"), e
	}
	t := terminal.NewTerminal(os.Stdin, "")
	oldState, _ := terminal.MakeRaw(0)
	defer terminal.Restore(0, oldState)
	p, e := t.ReadPassword("")
	return p, e
}

func enterPassword(name string) string {
	fmt.Printf("Enter password for %s: ", name)
	pass, e := readPassword()
	if e != nil {
		fmt.Println(e)
	}
	fmt.Printf("Retype password for %s: ", name)
	passAgain, e := readPassword()
	if e != nil {
		fmt.Println(e)
	}
	if pass != passAgain {
		fmt.Println("Error: the entered passwords do not match.")
		os.Exit(1)
	}
	return pass
}

func encryptPassword(pass, file string) {
	cmd := "echo \"" + pass + "\" | gpg -e " + getRecepientOpts() + " -o " + file + " --quiet --yes --compress-algo=none --no-encrypt-to"
	gpg := exec.Command("bash", "-c", cmd)
	if err := gpg.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func encryptMultiLine(file string) {
	scanner := bufio.NewScanner(os.Stdin)
	pass := []string{}
	for scanner.Scan() {
		pass = append(pass, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	encryptPassword(strings.Join(pass, "\n"), file)
}
