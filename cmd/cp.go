package cmd

import (
	"fmt"
	"os"
	"strings"

	"gitlab.com/JanMa/go-pass/pkg/git"
	"gitlab.com/JanMa/go-pass/pkg/store"
	"gitlab.com/JanMa/go-pass/pkg/store/entry"
	"gitlab.com/JanMa/go-pass/pkg/util"
)

func copyPasswords(src, dst string, force bool) (string, string) {
	dir, err := os.Stat(PasswordStore.Path + "/" + src)
	if !os.IsNotExist(err) && dir.IsDir() {
		entries, err := PasswordStore.FindEntries(src + "/.*")
		exitOnError(err)
		for _, e := range entries {
			copyPasswords(e.Name, strings.Replace(e.Name, src, dst, 1), force)
		}
		return PasswordStore.Path + "/" + src, PasswordStore.Path + "/" + dst
	}
	from, err := PasswordStore.FindEntry(src)
	exitOnError(err)
	to, err := PasswordStore.FindEntry(dst)
	if err == nil && !force &&
		!util.YesNo(fmt.Sprintf("%s already exists. Do you want to overwrite it?", dst)) {
		os.Exit(1)
	}
	to = entry.New(dst, PasswordStore.Path+"/"+dst+".gpg")
	err = from.Decrypt()
	exitOnError(err)
	v, err := from.Value()
	exitOnError(err)
	to.Insert(v)
	PasswordStore.InsertEntry(to)
	iD := PasswordStore.FindGpgID(to.Path)
	recv, err := store.ParseGpgID(iD)
	to.Encrypt(recv)
	git.AddFile(to.Path, fmt.Sprintf("Copy %s to %s.", src, dst))
	return from.Path, to.Path
}
