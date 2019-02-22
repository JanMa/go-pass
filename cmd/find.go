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

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find",
	Short: "List passwords that match pass-names",
	Run:   findPasswords,
}

func init() {
	rootCmd.AddCommand(findCmd)
}

func findPasswords(cmd *cobra.Command, args []string) {
	root := util.GetHomeDir() + "/.password-store"
	pattern := ""
	if len(args) > 0 {
		for _, a := range args {
			pattern += a
		}
	}
	// TODO: don't use external program
	lines := util.RunCommand("tree", "-C", "-l", "--noreport", "-P", pattern+"*", "--prune", "--matchdirs", "--ignore-case", root)
	fmt.Println("Search Terms:", pattern)
	for i := 1; i < len(lines); i++ {
		util.PrintLine(lines[i])
	}
}
