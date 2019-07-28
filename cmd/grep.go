package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/pkg/util"
)

func grepPasswords(cmd *cobra.Command, args []string) {
	r := regexp.MustCompile(args[0])
	all := PasswordStore.ShowAll()
	for _, e := range all {
		err := e.Decrypt()
		exitOnError(err)
		v, _ := e.Value()
		if r.MatchString(v) {
			fmt.Printf(util.BoldBlue+"%s:\n"+util.Reset, e.Name)
			for _, l := range strings.Split(v, "\n") {
				if r.MatchString(l) {
					fmt.Println(l)
				}
			}
		}
	}
}
