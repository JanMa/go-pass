package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/util"
)

// lsCmd represents the ls command
var (
	lsCmd = &cobra.Command{
		Use:                   "ls [pass-name]",
		Aliases:               []string{"list"},
		Short:                 "List passwords.",
		Args:                  cobra.MaximumNArgs(1),
		Run:                   listPasswords,
		DisableFlagsInUseLine: true,
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
		root += "/" + path
	}
	// TODO: don't use external program
	lines := util.RunCommand("tree", "-C", "-l", root, "-P", "*.gpg", "--noreport")
	fmt.Println(path)
	for i := 1; i < len(lines); i++ {
		util.PrintLine(lines[i])
	}
}
