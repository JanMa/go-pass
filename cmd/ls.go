package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/util"
)

func listPasswords(cmd *cobra.Command, args []string) {
	root := util.GetPasswordStore()
	path := "Password Store"
	if len(args) > 0 {
		path = args[0]
		root += "/" + path
	}
	// TODO: don't use external program
	lines := util.RunCommand("tree", "-C", "-l", root, "-P", "*.gpg", "--noreport")
	fmt.Println(path)
	for i := 1; i < len(lines); i++ {
		util.PrintLine(lines[i])
	}
}
