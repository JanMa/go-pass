package cmd

import (
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/pkg/git"
)

func gitCommand(cmd *cobra.Command, args []string) {
	exitOnError(git.RunCommand(args))
	if args[0] == "init" {
		git.AddFile(PasswordStore.Path, "Add current contents of password store.")
		attr := PasswordStore.Path + "/.gitattributes"
		os.Remove(attr)
		ioutil.WriteFile(attr, []byte("*.gpg diff=gpg\n"), 0666)
		git.AddFile(attr, "Configure git repository for gpg file diff.")
		git.RunCommand([]string{"config", "--local", "diff.gpg.binary", "true"})
		git.RunCommand([]string{"config", "--local", "diff.gpg.textconv", "gpg -d"})
	}
}
