package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/util"
	git "gopkg.in/src-d/go-git.v4"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		root := util.GetPasswordStore()
		_, e := git.PlainOpen(root)
		if e != nil {
			fmt.Println(e)
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(gitCmd)
}
