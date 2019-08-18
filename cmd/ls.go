package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/pkg/util"
)

func listPasswords(cmd *cobra.Command, args []string) {
	path := "Password Store"
	all := PasswordStore.ShowAll()
	if len(args) > 0 {
		path = args[0]
		all, _ = PasswordStore.FindEntries(args[0] + ".*")
	}
	names := []string{}
	for n := range all {
		names = append(names, n)
	}
	sort.Strings(names)
	fmt.Printf(util.BoldBlue+"%s:\n"+util.Reset, path)
	for _, n := range names {
		fmt.Println(all[n].Name)
	}
}
