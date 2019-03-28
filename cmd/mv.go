package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// mvCmd represents the mv command
var (
	mvCmd = &cobra.Command{
		Use:   "mv",
		Args:  cobra.ExactArgs(2),
		Short: "Renames or moves old-path to new-path, optionally forcefully, selectively reencrypting.",
		Run: func(cmd *cobra.Command, args []string) {
			s, _ := copyPasswords(args[0], args[1], ForceMv)
			if len(s) > 0 {
				os.RemoveAll(s)
			}
		},
	}

	ForceMv bool
)

func init() {
	rootCmd.AddCommand(mvCmd)

	mvCmd.Flags().BoolVarP(&ForceMv, "force", "f", false, "Forcefully copy password or directory.")
}
