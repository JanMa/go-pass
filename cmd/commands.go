package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "go-pass [subfolder | command]",
		Short: "go-pass is a pass clone written in Go",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			showPassword(cmd, args)
		},
		Example: "",
	}
	// versionCmd represents the version command
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Version:\t", Version)
			fmt.Println("Go version:\t", runtime.Version())
		},
	}
	// cpCmd represents the cp command
	cpCmd = &cobra.Command{
		Use:   "cp [--force, -f] old-path new-path",
		Args:  cobra.ExactArgs(2),
		Short: "Copies old-path to new-path, optionally forcefully, selectively reencrypting.",
		Run: func(cmd *cobra.Command, args []string) {
			copyPasswords(args[0], args[1], ForceCp)
		},
		Aliases:               []string{"copy"},
		DisableFlagsInUseLine: true,
	}
	// editCmd represents the edit command
	editCmd = &cobra.Command{
		Use:                   "edit pass-name",
		Short:                 "Insert a new password or edit an existing password using " + getEditor() + ".",
		Args:                  cobra.ExactArgs(1),
		Run:                   editPassword,
		DisableFlagsInUseLine: true,
	}
	// findCmd represents the find command
	findCmd = &cobra.Command{
		Use:                   "find pass-names...",
		Short:                 "List passwords that match pass-names",
		Args:                  cobra.MinimumNArgs(1),
		Run:                   findPasswords,
		Aliases:               []string{"search"},
		DisableFlagsInUseLine: true,
	}
	// generateCmd represents the generate command
	generateCmd = &cobra.Command{
		Use:   "generate [--no-symbols,-n] [--clip,-c] [--qrcode,-q] [--in-place,-i | --force,-f] pass-name [pass-length]",
		Args:  cobra.RangeArgs(1, 2),
		Short: "Generate a new password of pass-length (or 25 if unspecified) with optionally no symbols.",
		Long: `Generate a new password of pass-length (or 25 if unspecified) with optionally no symbols.
Optionally put it on the clipboard and clear board after 45 seconds.
Prompt before overwriting existing password unless forced.
Optionally replace only the first line of an existing file with a new password.`,
		Run:                   generatePassword,
		DisableFlagsInUseLine: true,
	}
	// gitCmd represents the git command
	gitCmd = &cobra.Command{
		Use:                "git",
		Short:              "If the password store is a git repository, execute a git command specified by git-command-args.",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		Run:                git,
	}
	// grepCmd represents the grep command
	grepCmd = &cobra.Command{
		Use:                   "grep [GREPOPTIONS] search-string",
		Args:                  cobra.MinimumNArgs(1),
		Short:                 "Search for password files containing search-string when decrypted.",
		Run:                   grepPasswords,
		DisableFlagsInUseLine: true,
	}
	// initCmd represents the init command
	initCmd = &cobra.Command{
		Use:   "init [--path=subfolder,-p subfolder] gpg-id...",
		Args:  cobra.MinimumNArgs(1),
		Short: "Initialize new password storage and use gpg-id for encryption.",
		Long: `Initialize new password storage and use gpg-id for encryption.
Selectively reencrypt existing passwords using new gpg-id.`,
		Run:                   initPasswordStore,
		DisableFlagsInUseLine: true,
	}
	// insertCmd represents the insert command
	insertCmd = &cobra.Command{
		Use:   "insert [--echo,-e | --multiline,-m] [--force,-f] pass-name",
		Args:  cobra.ExactArgs(1),
		Short: "Insert new password.",
		Long: `Insert new password. Optionally, echo the password back to the console
during entry. Or, optionally, the entry may be multiline. Prompt before
overwriting existing password unless forced.`,
		Run:                   insertPassword,
		Aliases:               []string{"add"},
		DisableFlagsInUseLine: true,
	}
	// lsCmd represents the ls command
	lsCmd = &cobra.Command{
		Use:                   "ls [pass-name]",
		Aliases:               []string{"list"},
		Short:                 "List passwords.",
		Args:                  cobra.MaximumNArgs(1),
		Run:                   listPasswords,
		DisableFlagsInUseLine: true,
	}
	// mvCmd represents the mv command
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
	// rmCmd represents the rm command
	rmCmd = &cobra.Command{
		Use:                   "rm [--recursive,-r] [--force,-f] pass-name",
		Args:                  cobra.ExactArgs(1),
		Short:                 "Remove existing password or directory, optionally forcefully.",
		Run:                   rmPassword,
		Aliases:               []string{"delete", "remove"},
		DisableFlagsInUseLine: true,
	}
	// showCmd represents the show command
	showCmd = &cobra.Command{
		Use:                   "show [--clip[=line-number],-c[=line-number]] [--qrcode[=line-number],-q[=line-number]] [pass-name]",
		Short:                 "Show existing password and optionally put it on the clipboard.",
		Args:                  cobra.MaximumNArgs(1),
		Run:                   showPassword,
		Aliases:               []string{"ls", "list"},
		DisableFlagsInUseLine: true,
	}
)

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
