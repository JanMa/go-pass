package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/util"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		root := util.GetPasswordStore()
		gitStatus := exec.Command("git", "-C", root, "status")
		if o, e := gitStatus.CombinedOutput(); e != nil {
			fmt.Println(o)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(gitCmd)
}
