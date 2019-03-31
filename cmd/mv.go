package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// mvCmd represents the mv command
var (
	mvCmd = &cobra.Command{
		Use:   "mv [--force,-f] old-path new-path",
		Args:  cobra.ExactArgs(2),
		Short: "Renames or moves old-path to new-path, optionally forcefully, selectively reencrypting.",
		Run: func(cmd *cobra.Command, args []string) {
			s, _ := copyPasswords(args[0], args[1], ForceMv)
			if len(s) > 0 {
				os.RemoveAll(s)
				gitAddFile(strings.TrimRight(s, "/"), fmt.Sprintf("Remove %s from store.", args[0]))
			}
		},
		Aliases:               []string{"rename"},
		DisableFlagsInUseLine: true,
	}

	ForceMv bool
)

func init() {
	rootCmd.AddCommand(mvCmd)

	mvCmd.Flags().BoolVarP(&ForceMv, "force", "f", false, "Forcefully copy password or directory.")
}
