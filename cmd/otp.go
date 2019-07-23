package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/pkg/util"
)

func genOtpCode(cmd *cobra.Command, args []string) {
	root := util.GetPasswordStore() + "/" + args[0] + ".gpg"

	if _, e := os.Stat(root); os.IsNotExist(e) {
		fmt.Printf("Error: %s is not in the password store.\n", args[0])
		os.Exit(1)
	}
	lines := util.RunCommand("gpg", "-dq", root)
	u, e := findOtpUrl(lines)
	exitOnError(e)
	k, e := otp.NewKeyFromURL(u)
	exitOnError(e)
	c, e := totp.GenerateCode(k.Secret(), time.Now())
	exitOnError(e)
	if OtpClip {
		exitOnError(clipboard.WriteAll(c))
	} else {
		fmt.Println(c)
	}
}

func findOtpUrl(lines []string) (string, error) {
	for _, l := range lines {
		if strings.Contains(l, "otpauth://totp") {
			return l, nil
		}
	}
	return "", fmt.Errorf("OTP secret not found.")
}
