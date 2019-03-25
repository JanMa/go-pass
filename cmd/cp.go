package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/copy"
	"gitlab.com/JanMa/go-pass/util"
)

// cpCmd represents the cp command
var (
	cpCmd = &cobra.Command{
		Use:   "cp",
		Args:  cobra.ExactArgs(2),
		Short: "Copies old-path to new-path, optionally forcefully, selectively reencrypting.",
		Run:   copyPassword,
	}

	ForceCp bool
)

func init() {
	rootCmd.AddCommand(cpCmd)

	cpCmd.Flags().BoolVarP(&ForceCp, "force", "f", false, "Forcefully copy password or directory.")
}

func copyPassword(cmd *cobra.Command, args []string) {
	fromPath := util.GetPasswordStore() + "/" + args[0]
	toPath := util.GetPasswordStore() + "/" + args[1]
	if f, e := os.Stat(fromPath); !os.IsNotExist(e) && f.IsDir() {
		if f, e := os.Stat(toPath); !os.IsNotExist(e) && f.IsDir() && !ForceCp &&
			!util.YesNo(fmt.Sprintf("%s already exists. Do you want to overwrite it?", toPath)) {
			os.Exit(1)
		}
		fmt.Println(fromPath)
		copy.Copy(fromPath, toPath)
	} else if f, e := os.Stat(fromPath + ".gpg"); !os.IsNotExist(e) && !f.IsDir() {
		if strings.HasSuffix(toPath, "/") {
			toPath += filepath.Base(fromPath)
		}
		if f, e := os.Stat(toPath + ".gpg"); !os.IsNotExist(e) && !f.IsDir() && !ForceCp &&
			!util.YesNo(fmt.Sprintf("%s already exists. Do you want to overwrite it?", toPath+".gpg")) {
			os.Exit(1)
		}
		fmt.Println(fromPath + ".gpg")
		copy.Copy(fromPath+".gpg", toPath+".gpg")
		fmt.Printf("'%s.gpg' -> '%s.gpg'\n", fromPath, toPath)
	} else {
		fmt.Printf("Error: %s is not in the password store.\n", args[0])
		os.Exit(1)
	}
}
