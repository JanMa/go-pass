package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/pkg/git"
	"gitlab.com/JanMa/go-pass/pkg/gpg"
	"gitlab.com/JanMa/go-pass/pkg/store"
)

func initPasswordStore(cmd *cobra.Command, args []string) {
	gpgKeys, err := store.GetGpgKeys(args)
	exitOnError(err)
	root := PasswordStore.Path
	if len(Subdir) > 0 {
		root += "/" + strings.TrimRight(Subdir, "/")
	}
	gpgID := root + "/.gpg-id"
	if _, e := os.Stat(gpgID); os.IsExist(e) {
		exitOnError(os.Remove(gpgID))
	}
	exitOnError(ioutil.WriteFile(gpgID, []byte(strings.Join(args, "\n")+"\n"), 0644))
	fmt.Printf("Password store initialized for %s\n",
		strings.Trim(strings.Join(args, ", "), "\n"))
	git.AddFile(gpgID, fmt.Sprintf("Set GPG id to %s.",
		strings.Trim(strings.Join(args, ", "), "\n")))
	if len(Subdir) > 0 {
		exitOnError(reEncryptEntries(Subdir+"/", gpgKeys))
	} else {
		exitOnError(reEncryptEntries("", gpgKeys))
	}
	git.AddFile(root, fmt.Sprintf("Reencrypt password store using new GPG id %s.",
		strings.Trim(strings.Join(args, ", "), "\n")))
}

func reEncryptEntries(path string, keys []string) error {
	if len(PasswordStore.ShowAll()) == 0 {
		return nil
	}
	entries, err := PasswordStore.FindEntries(path + ".*")
	names := store.SortEntries(entries)
	if err != nil {
		return err
	}
	for _, n := range names {
		curKeys, err := gpg.GetKeys(entries[n].Path)
		if err != nil {
			return err
		}
		if !matchKeys(curKeys, keys) {
			err := entries[n].Decrypt()
			if err != nil {
				return err
			}
			fmt.Println("reencrypting", n)
			err = entries[n].Encrypt(keys)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func matchKeys(input string, keys []string) bool {
	for _, k := range keys {
		if ret, _ := regexp.MatchString(k, input); !ret {
			return false
		}
	}
	return true
}
