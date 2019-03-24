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

// showCmd represents the show command
var (
	showCmd = &cobra.Command{
		Use:   "show",
		Short: "Show existing password and optionally put it on the clipboard.",
		Args:  cobra.MaximumNArgs(1),
		Run:   showPassword,
	}

	Copy   int
	QRCode int
)

func init() {
	rootCmd.AddCommand(showCmd)

	showCmd.Flags().IntVarP(&Copy, "clip", "c", 0, "Copy password to clipboard")
	showCmd.Flags().Lookup("clip").NoOptDefVal = "1"
	showCmd.Flags().IntVarP(&QRCode, "qrcode", "q", 0, "Display output as QR code")
	showCmd.Flags().Lookup("qrcode").NoOptDefVal = "1"
}

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
		// TODO: don't use external program
		lines := util.RunCommand("gpg", "-dq", root)
		if Copy > 0 {
			if Copy > len(lines) {
				fmt.Printf("There is no password to put on the clipboard at line %d.\n", Copy)
				os.Exit(1)
			}
			if err := clipboard.WriteAll(lines[Copy-1]); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if QRCode > 0 {
			if QRCode > len(lines) {
				fmt.Printf("There is no password to put on the clipboard at line %d.\n", QRCode)
				os.Exit(1)
			}
			qr, err := qrcode.New(lines[QRCode-1], qrcode.Low)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Print(qr.ToSmallString(false))
		} else {
			fmt.Print(strings.Join(lines, "\n"))
		}
	}
}
