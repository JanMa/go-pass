package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/util"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:                "git",
	Short:              "If the password store is a git repository, execute a git command specified by git-command-args.",
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		root := util.GetPasswordStore()
		if f, e := os.Stat(root + "/.git"); (os.IsNotExist(e) || !f.IsDir()) && args[0] != "init" {
			fmt.Printf("Error: the password store is not a git repository. Try \"go-pass git init\".\n")
			os.Exit(1)
		}
		git := exec.Command("git", "-C", root)
		for _, a := range args {
			git.Args = append(git.Args, a)
		}
		o, _ := git.CombinedOutput()
		fmt.Print(string(o))
		if args[0] == "init" {
			gitAddFile(root, "Add current contents of password store.")
			attr := root + "/.gitattributes"
			os.Remove(attr)
			ioutil.WriteFile(attr, []byte("*.gpg diff=gpg\n"), 0666)
			gitAddFile(attr, "Configure git repository for gpg file diff.")
			exec.Command("git", "-C", root, "config", "--local", "diff.gpg.binary", "true").Run()
			exec.Command("git", "-C", root, "config", "--local", "diff.gpg.textconv", "gpg -d").Run()
		}
	},
}

func gitAddFile(path, msg string) error {
	r := util.GetPasswordStore()
	if f, e := os.Stat(r + "/.git"); os.IsNotExist(e) || !f.IsDir() {
		return nil
	}
	o, e := exec.Command("git", "-C", r, "add", path).CombinedOutput()
	fmt.Print(string(o))
	if e != nil {
		return e
	}
	s, _ := exec.Command("git", "-C", r, "status", "--porcelain", path).CombinedOutput()
	if len(msg) > 0 && len(s) > 0 {
		sign := "--no-gpg-sign"
		if c, _ := exec.Command("git", "-C", r, "config", "--bool", "--get", "pass.signcommits").CombinedOutput(); string(c) == "true" {
			sign = "-S"
		}
		o, e := exec.Command("git", "-C", r, "commit", sign, "-m", msg).CombinedOutput()
		fmt.Print(string(o))
		if e != nil {
			return e
		}
	}
	return nil
}
