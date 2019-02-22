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
	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/util"
)

// lsCmd represents the ls command
var (
	lsCmd = &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List passwords.",
		Args:    cobra.MinimumNArgs(0),
		Run:     listPasswords,
	}
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

func listPasswords(cmd *cobra.Command, args []string) {
	root := util.GetPasswordStore()
	path := "Password Store"
	if len(args) > 0 {
		path = args[0]
		for _, a := range args {
			root += "/" + a
		}
	}
	// TODO: don't use external program
	lines := util.RunCommand("tree", root, "-P", "*.gpg", "--noreport")
	fmt.Println(path)
	for i := 1; i < len(lines); i++ {
		util.PrintLine(lines[i])
	}
}
