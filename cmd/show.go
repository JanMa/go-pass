package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/util"
)

func showPassword(cmd *cobra.Command, args []string) {
	root := util.GetPasswordStore()
	if len(args) > 0 {
		root += "/" + args[0]
	}
	if f, e := os.Stat(root); !os.IsNotExist(e) && f.IsDir() {
		listPasswords(cmd, args)
	} else {
		root += ".gpg"
		if _, e := os.Stat(root); os.IsNotExist(e) {
			fmt.Printf("Error: %s is not in the password store.\n", args[0])
			os.Exit(1)
		}
		lines := util.RunCommand("gpg", "-dq", root)
		if Copy > 0 {
			if Copy > len(lines) {
				fmt.Printf("There is no password to put on the clipboard at line %d.\n", Copy)
				os.Exit(1)
			}
			exitOnError(clipboard.WriteAll(lines[Copy-1]))
		} else if QRCode > 0 {
			if QRCode > len(lines) {
				fmt.Printf("There is no password to put on the clipboard at line %d.\n", QRCode)
				os.Exit(1)
			}
			qr, err := qrcode.New(lines[QRCode-1], qrcode.Low)
			exitOnError(err)
			fmt.Print(qr.ToSmallString(false))
		} else {
			fmt.Print(strings.Join(lines, "\n"))
		}
	}
}
