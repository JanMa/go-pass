package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/pkg/git"
	"gitlab.com/JanMa/go-pass/pkg/util"
)

func rmPassword(cmd *cobra.Command, args []string) {
	pattern := args[0]
	dir, err := os.Stat(PasswordStore.Path + "/" + args[0])
	if !os.IsNotExist(err) && dir.IsDir() {
		pattern = args[0] + "/.*"
	}
	result, err := PasswordStore.FindEntries(pattern)
	if err != nil {
		fmt.Println("found no matching entries for", args[0])
		os.Exit(1)
	}
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
