package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func findPasswords(cmd *cobra.Command, args []string) {
	result, err := PasswordStore.FindEntries(args[0])
	exitOnError(err)
	fmt.Println("Results:")
	for _, entry := range result {
		fmt.Println(entry.Name)
	}
}
