package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/cobra"
)

func genOtpCode(cmd *cobra.Command, args []string) {
	e, err := PasswordStore.FindEntry(args[0])
	exitOnError(err)
	exitOnError(e.Decrypt())
	v, _ := e.Value()
	lines := strings.Split(v, "\n")
	u, err := findOtpURL(lines)
	exitOnError(err)
	k, err := otp.NewKeyFromURL(u)
	exitOnError(err)
	c, err := totp.GenerateCode(k.Secret(), time.Now())
	exitOnError(err)
	if OtpClip {
		exitOnError(clipboard.WriteAll(c))
	} else {
		fmt.Println(c)
	}
}

func findOtpURL(lines []string) (string, error) {
	for _, l := range lines {
		if strings.Contains(l, "otpauth://totp") {
			return l, nil
		}
	}
	return "", fmt.Errorf("otp secret not found")
}
