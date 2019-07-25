package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/pkg/git"
	"gitlab.com/JanMa/go-pass/pkg/util"
)

func rmPassword(cmd *cobra.Command, args []string) {
	result, err := PasswordStore.FindEntries(args[0])
	exitOnError(err)
	fmt.Println("The following entries will be deleted:")
	for _, entry := range result {
		fmt.Println("-", entry.Name)
	}
	fmt.Println()
	if !ForceRm && !util.YesNo(fmt.Sprintf("Are you sure you would like to delete them?")) {
		os.Exit(1)
	}

	for _, entry := range result {
		err = entry.Delete()
		exitOnError(err)
		git.AddFile(entry.Path, fmt.Sprintf("Remove %s from store.", entry.Name))
	}
}
