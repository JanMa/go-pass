package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
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

	// CharacterSet set containing all types of characters
	CharacterSet = LowerLetters + UpperLetters + Digits + Symbols
	// CharacterSetNoSymbols set containing only letters and digits
	CharacterSetNoSymbols = LowerLetters + UpperLetters + Digits

	// BoldBlue print text in bold blue
	BoldBlue = "\033[1m\033[34m"
	// Reset reset text formating
	Reset = "\033[0m"
)

// YesNo simple Yes or No dialogue
func YesNo(msg string) bool {
	t := terminal.NewTerminal(os.Stdin, "")
	oldState, _ := terminal.MakeRaw(0)
	defer terminal.Restore(0, oldState)
	fmt.Printf("%s [y/N] ", msg)
	i, e := t.ReadLine()
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	i = strings.ToLower(i)
	if i == "y" || i == "yes" {
		return true
	}
	return false
}

// RandomString returns a random string of of the given length
func RandomString(n int, symbols bool) (string, error) {
	letter := []rune(CharacterSet)
	if !symbols {
		letter = []rune(CharacterSetNoSymbols)
	}

	b := make([]rune, n)
	for i := 0; i < n; i++ {
		c, err := rand.Int(rand.Reader, big.NewInt(int64(len(letter))))
		if err != nil {
			return "", err
		}
		b[i] = letter[c.Int64()]
	}
	return string(b), nil
}

// TmpFile generates the path to a new temporary file
func TmpFile() (string, error) {
	os.Mkdir(os.TempDir()+"/go-pass", 0700)
	r, e := RandomString(8, false)
	return os.TempDir() + "/go-pass/" + r, e
}
