package cmd

import (
	"fmt"
	"io/ioutil"
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
		Run: func(cmd *cobra.Command, args []string) {
			copyPasswords(args[0], args[1], ForceCp)
		},
	}

	ForceCp bool
)

func init() {
	rootCmd.AddCommand(cpCmd)

	cpCmd.Flags().BoolVarP(&ForceCp, "force", "f", false, "Forcefully copy password or directory.")
}

func copyPasswords(src, dst string, force bool) {
	fromPath := util.GetPasswordStore() + "/" + src
	toPath := util.GetPasswordStore() + "/" + dst
	// src exists and is a directory
	if f, e := os.Stat(fromPath); !os.IsNotExist(e) && f.IsDir() {
		// dst exists and is a directory
		if f, e := os.Stat(toPath); !os.IsNotExist(e) && f.IsDir() && !force &&
			!util.YesNo(fmt.Sprintf("%s already exists. Do you want to overwrite it?", toPath)) {
			os.Exit(1)
		}
		fmt.Println(fromPath)
		copy.Copy(fromPath, toPath)
		recv, e := getRecipientsFromGpgID(findGpgID(toPath))
		if e != nil {
			fmt.Println(1)
			os.Exit(1)
		}
		// walk dst directory
		reEncryptDir(toPath, recv)
		// src exists and is not a directory
	} else if f, e := os.Stat(fromPath + ".gpg"); !os.IsNotExist(e) && !f.IsDir() {
		// dst has a slash as suffix indicating it is a directory
		if strings.HasSuffix(toPath, "/") {
			toPath += filepath.Base(fromPath)
		}
		// dst exists and is a file and no forcefuly overwriting
		if f, e := os.Stat(toPath + ".gpg"); !os.IsNotExist(e) && !f.IsDir() && !force &&
			!util.YesNo(fmt.Sprintf("%s already exists. Do you want to overwrite it?", toPath+".gpg")) {
			os.Exit(1)
		}
		fmt.Println(fromPath + ".gpg")
		copy.Copy(fromPath+".gpg", toPath+".gpg")
		recv, e := getRecipientsFromGpgID(findGpgID(toPath))
		if e != nil {
			fmt.Println(1)
			os.Exit(1)
		}
		reEncryptFile(toPath+".gpg", recv)
	} else {
		fmt.Printf("Error: %s is not in the password store.\n", src)
		os.Exit(1)
	}
}

func findGpgID(path string) string {
	dirs := strings.Split(strings.Trim(path, "/"), "/")
	root := strings.Split(strings.Trim(util.GetPasswordStore(), "/"), "/")
	if l := len(dirs); filepath.Ext(dirs[l-1]) == ".gpg" {
		dirs = dirs[:(l - 1)]
	}
	for i := len(dirs); i >= len(root); i-- {
		p := "/" + strings.Join(dirs[:i], "/") + "/.gpg-id"
		if f, e := os.Stat(p); !os.IsNotExist(e) && !f.IsDir() {
			return p
		}
	}
	return util.GetPasswordStore() + "/.gpg-id"
}

func getRecipientsFromGpgID(path string) ([]string, error) {
	if _, e := os.Stat(path); os.IsNotExist(e) {
		return nil, e
	}
	gpgID, e := ioutil.ReadFile(path)
	if e != nil {
		return nil, e
	}
	return strings.Split(strings.Trim(string(gpgID), "\n"), "\n"), nil
}
