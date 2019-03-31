package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Version:\t", Version)
			fmt.Println("Go version:\t", runtime.Version())
		},
	}

	Version string
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
