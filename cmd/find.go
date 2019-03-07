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
	root := util.GetPasswordStore()
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
