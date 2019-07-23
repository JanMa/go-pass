package main

import (
	"fmt"
	"os"

	"gitlab.com/JanMa/go-pass/pkg/git"
	"gitlab.com/JanMa/go-pass/pkg/store"
	"gitlab.com/JanMa/go-pass/pkg/store/entry"
)

func main() {

	s, err := store.GetPasswordStore()
	pStore, err := store.New(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = pStore.Fill()

	name := "test/Test"
	newEntry := entry.New(name, s+"/"+name+".gpg")
	newEntry.Insert("A test!")
	gpgID := pStore.FindGpgID(newEntry.Path)
	recp, err := store.ParseGpgID(gpgID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = newEntry.Encrypt(recp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pStore.InsertEntry(newEntry)

	az, err := pStore.FindEntry(name)
	err = az.Decrypt()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = az.Show()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	g, err := git.AddFile(az.Path, "")
	fmt.Println(g)
	status, err := git.RunCommand([]string{"status"})
	fmt.Println(status)
}
