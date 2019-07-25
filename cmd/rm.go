package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/pkg/git"
	"gitlab.com/JanMa/go-pass/pkg/util"
)

func rmPassword(cmd *cobra.Command, args []string) {
	pattern := args[0]
	dir := PasswordStore.Path + "/" + args[0]
	stat, err := os.Stat(dir)
	if !os.IsNotExist(err) && stat.IsDir() {
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

	// Ensure empty directory gets deleted
	if empty, _ := isEmpty(dir); empty {
		fmt.Println("Remove:", dir)
		err = os.RemoveAll(dir)
		git.AddFile(dir, fmt.Sprintf("Remove %s from store.", args[0]))
	}
}

// https://stackoverflow.com/a/30708914
func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}
