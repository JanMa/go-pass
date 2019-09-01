package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/pkg/store"
	"gitlab.com/JanMa/go-pass/pkg/util"
)

func findPasswords(cmd *cobra.Command, args []string) {
	entries, err := PasswordStore.FindEntries(args[0])
	exitOnError(err)
	names := store.SortEntries(entries)
	fmt.Printf(util.BoldBlue + "Results:\n" + util.Reset)
	for _, n := range names {
		fmt.Println(n)
	}
}
