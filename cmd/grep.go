package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/util"
)

func grepPasswords(cmd *cobra.Command, args []string) {
	grepArgs := strings.Join(args, " ")
	err := filepath.Walk(util.GetPasswordStore(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if filepath.Ext(path) == ".gpg" {
			p := strings.Join(util.RunCommand("gpg", "-dq", path), "\n")
			grep := exec.Command("grep", "--color=always", grepArgs)
			stdin, err := grep.StdinPipe()
			exitOnError(err)
			go func() {
				defer stdin.Close()
				io.WriteString(stdin, p)
			}()
			out, _ := grep.Output()
			if len(out) > 0 {
				rel, _ := filepath.Rel(util.GetPasswordStore(), path)
				util.PrintLine(rel)
				fmt.Print(string(out))
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}
