package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/util"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:                   "find pass-names...",
	Short:                 "List passwords that match pass-names",
	Args:                  cobra.MinimumNArgs(1),
	Run:                   findPasswords,
	Aliases:               []string{"search"},
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(findCmd)
}

func findPasswords(cmd *cobra.Command, args []string) {
	root := util.GetPasswordStore()
	pattern := "*"
	for _, a := range args {
		pattern += a + "*|*"
	}
	lines := util.RunCommand("tree", "-C", "-l", "--noreport", "-P", strings.TrimSuffix(pattern, "|*"), "--prune", "--matchdirs", "--ignore-case", root)
	fmt.Println("Search Terms:", strings.Join(args, " "))
	for _, l := range lines {
		util.PrintLine(l)
	}
}
