package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/pkg/git"
	"gitlab.com/JanMa/go-pass/pkg/store"
	"gitlab.com/JanMa/go-pass/pkg/store/entry"
	"gitlab.com/JanMa/go-pass/pkg/util"
)

var (
	verb = map[bool]string{
		true:  "Replace",
		false: "Add",
	}
)

func generatePassword(cmd *cobra.Command, args []string) {
	length := 25
	// args parsing
	if len(args) > 1 {
		l, e := strconv.Atoi(args[1])
		if e != nil {
			fmt.Println("Error: pass length", args[1], "must be a number.")
			os.Exit(1)
		}
		if l == 0 {
			fmt.Println("Error: pass-length must be greater than zero.")
			os.Exit(1)
		}
		length = l
	}
	pass, err := util.RandomString(length, !NoSymbols)
	exitOnError(err)
	//find existing entry
	e, err := PasswordStore.FindEntry(args[0])
	if err == nil {
		if !ForceGen &&
			!util.YesNo(fmt.Sprintf("An entry already exists for %s. Overwrite it?", args[0])) {
			os.Exit(1)
		}
		err = e.Decrypt()
		exitOnError(err)
		if !InPlace {
			e.Insert(pass + "\n")
		} else {
			oldVal, _ := e.Value()
			l := strings.Split(oldVal, "\n")
			l[0] = pass
			e.Insert(strings.Join(l, "\n"))
		}
	} else {
		e = entry.New(args[0], PasswordStore.Path+"/"+args[0]+".gpg")
		PasswordStore.InsertEntry(e)
		e.Insert(pass + "\n")
	}

	iD := PasswordStore.FindGpgID(e.Path)
	recv, _ := store.ParseGpgID(iD)
	exitOnError(e.Encrypt(recv))

	git.AddFile(e.Path, fmt.Sprintf("%s generated password for %s.", verb[InPlace], args[0]))

	if Clip {
		exitOnError(clipboard.WriteAll(pass))
	} else if GenQRCode {
		qr, err := qrcode.New(pass, qrcode.Low)
		exitOnError(err)
		fmt.Print(qr.ToSmallString(false))
	} else {
		fmt.Printf("The generated password for %s is:\n%s\n", args[0], pass)
	}
}
