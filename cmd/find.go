package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/pkg/util"
)

func findPasswords(cmd *cobra.Command, args []string) {
	root := util.GetPasswordStore()
	pattern := "*"
	for _, a := range args {
		pattern += a + "*|*"
	}
	lines := util.RunCommand("tree", "-C", "-l", "--noreport", "-P", strings.TrimSuffix(pattern, "|*"), "--prune", "--matchdirs", "--ignore-case", root)
	fmt.Println("Search Terms:", strings.Join(args, " "))
	for _, l := range lines[1:] {
		util.PrintLine(l)
	}
}
