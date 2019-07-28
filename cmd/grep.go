package cmd

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/pkg/util"
)

func grepPasswords(cmd *cobra.Command, args []string) {
	r := regexp.MustCompile(args[0])
	all := PasswordStore.ShowAll()
	names := []string{}
	for n := range all {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		err := all[n].Decrypt()
		exitOnError(err)
		v, _ := all[n].Value()
		if r.MatchString(v) {
			fmt.Printf(util.BoldBlue+"%s:\n"+util.Reset, all[n].Name)
			for _, l := range strings.Split(v, "\n") {
				if r.MatchString(l) {
					fmt.Println(l)
				}
			}
		}
	}
}
