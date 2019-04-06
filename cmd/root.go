package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-pass [subfolder | command]",
	Short: "go-pass is a pass clone written in Go",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		showPassword(cmd, args)
	},
	Example: "",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	exitOnError(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(cpCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(findCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(gitCmd)
	rootCmd.AddCommand(grepCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(insertCmd)
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(mvCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(versionCmd)
}
