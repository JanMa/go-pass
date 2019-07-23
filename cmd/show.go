package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

func showPassword(cmd *cobra.Command, args []string) {
	e, err := PasswordStore.FindEntry(args[0])
	if err != nil {
		listPasswords(cmd, args)
		return
	}
	err = e.Decrypt()
	exitOnError(err)
	l, err := e.Value()
	lines := strings.Split(l, "\n")
	if Copy > 0 {
		if Copy > len(lines) {
			fmt.Printf("There is no password to put on the clipboard at line %d.\n", Copy)
			os.Exit(1)
		}
		exitOnError(clipboard.WriteAll(lines[Copy-1]))
		return
	}
	if QRCode > 0 {
		if QRCode > len(lines) {
			fmt.Printf("There is no password to put on the clipboard at line %d.\n", QRCode)
			os.Exit(1)
		}
		qr, err := qrcode.New(lines[QRCode-1], qrcode.Low)
		exitOnError(err)
		fmt.Print(qr.ToSmallString(false))
		return
	}
	fmt.Print(strings.Join(lines, "\n"))
}
