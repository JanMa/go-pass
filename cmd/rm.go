package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/util"
)

func rmPassword(cmd *cobra.Command, args []string) {
	passDir := util.GetPasswordStore() + "/" + args[0]
	passFile := passDir + ".gpg"
	_, eF := os.Stat(passFile)
	fD, eD := os.Stat(passDir)
	if !os.IsNotExist(eF) && !os.IsNotExist(eD) && fD.IsDir() && args[0][len(args[0])-1] == '/' || os.IsNotExist(eF) {
		passFile = strings.TrimRight(passDir, "/")
	}

	if _, e := os.Stat(passFile); os.IsNotExist(e) {
		fmt.Printf("Error: %s is not in the password store.\n", args[0])
		os.Exit(1)
	}

	if !ForceRm && !util.YesNo(fmt.Sprintf("Are you sure you would like to delete %s?", args[0])) {
		os.Exit(1)
	}

	if e := func(p string, r bool) error {
		if r {
			return os.RemoveAll(p)
		}
		return os.Remove(p)
	}(passFile, RecurseRm); e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	gitAddFile(strings.TrimRight(passFile, "/"), fmt.Sprintf("Remove %s from store.", args[0]))
}
