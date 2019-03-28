package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/util"
)

// generateCmd represents the generate command
var (
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

	NoSymbols bool
	Clip      bool
	InPlace   bool
	ForceGen  bool
	GenQRCode bool
)

const (
	// LowerLetters is the list of lowercase letters.
	LowerLetters = "abcdefghijklmnopqrstuvwxyz"

	// UpperLetters is the list of uppercase letters.
	UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Digits is the list of permitted digits.
	Digits = "0123456789"

	// Symbols is the list of symbols.
	Symbols = "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"

	CharacterSet          = LowerLetters + UpperLetters + Digits + Symbols
	CharacterSetNoSymbols = LowerLetters + UpperLetters + Digits
)

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().BoolVarP(&NoSymbols, "no-symbols", "n", false, "Generate password with no symbols.")
	generateCmd.Flags().BoolVarP(&Clip, "clip", "c", false, "Put generated password on the clipboard.")
	generateCmd.Flags().BoolVarP(&InPlace, "in-place", "i", false, "Replace only the first line of an existing file with a new password.")
	generateCmd.Flags().BoolVarP(&ForceGen, "force", "f", false, "Forcefully overwrite existing password.")
	generateCmd.Flags().BoolVarP(&GenQRCode, "qrcode", "q", false, "Display output as QR code.")

}
func generatePassword(cmd *cobra.Command, args []string) {
	length := 25
	if len(args) > 1 {
		l, e := strconv.Atoi(args[1])
		if e != nil {
			fmt.Println("Error: pass length", args[1], "must be a number.")
			os.Exit(1)
		}
		length = l
	}
	if length == 0 {
		fmt.Println("Error: pass-length must be greater than zero.")
		os.Exit(1)
	}
	root := util.GetPasswordStore() + "/" + args[0] + ".gpg"
	if f, e := os.Stat(root); !os.IsNotExist(e) && !f.IsDir() {
		overwrite := func() {
			if err := os.Remove(root); err != nil {
				fmt.Println(err)
			}
		}
		if ForceGen && !InPlace || !ForceGen && !InPlace &&
			util.YesNo(fmt.Sprintf("An entry already exists for %s. Overwrite it?", args[0])) {
			overwrite()
		} else if !InPlace {
			os.Exit(1)
		}
	}
	pass := randomString(length, !NoSymbols)
	if _, e := os.Stat(root); !os.IsNotExist(e) && InPlace {
		lines := util.RunCommand("gpg", "-dq", root)
		lines[0] = pass
		os.Remove(root)
		encryptPassword(strings.Join(lines, "\n"), root)

	} else {
		encryptPassword(pass+"\n", root)
	}

	if Clip {
		if err := clipboard.WriteAll(pass); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if GenQRCode {
		qr, err := qrcode.New(pass, qrcode.Low)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Print(qr.ToSmallString(false))
	} else {
		fmt.Printf("The generated password for %s is:\n%s\n", args[0], pass)
	}
}

func randomString(length int, symbols bool) string {
	pass := ""
	set := CharacterSet
	if !symbols {
		set = CharacterSetNoSymbols
	}
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(set))))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		pass += string(set[n.Int64()])
	}
	return pass
}
