package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/pkg/git"
	"gitlab.com/JanMa/go-pass/util"
)

func gitCommand(cmd *cobra.Command, args []string) {
	root := util.GetPasswordStore()
	if f, e := os.Stat(root + "/.git"); (os.IsNotExist(e) || !f.IsDir()) && args[0] != "init" {
		fmt.Printf("Error: the password store is not a git repository. Try \"go-pass git init\".\n")
		os.Exit(1)
	}
	run := exec.Command("git", "-C", root)
	for _, a := range args {
		run.Args = append(run.Args, a)
	}
	o, _ := run.CombinedOutput()
	fmt.Print(string(o))
	if args[0] == "init" {
		git.AddFile(root, "Add current contents of password store.")
		attr := root + "/.gitattributes"
		os.Remove(attr)
		ioutil.WriteFile(attr, []byte("*.gpg diff=gpg\n"), 0666)
		git.AddFile(attr, "Configure git repository for gpg file diff.")
		exec.Command("git", "-C", root, "config", "--local", "diff.gpg.binary", "true").Run()
		exec.Command("git", "-C", root, "config", "--local", "diff.gpg.textconv", "gpg -d").Run()
	}
}
