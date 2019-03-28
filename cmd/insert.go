package cmd

import (
	"bufio"
	"fmt"
	"gitlab.com/JanMa/go-pass/util"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// insertCmd represents the insert command
var (
	insertCmd = &cobra.Command{
		Use:   "insert [--echo,-e | --multiline,-m] [--force,-f] pass-name",
		Args:  cobra.ExactArgs(1),
		Short: "Insert new password.",
		Long: `Insert new password. Optionally, echo the password back to the console
during entry. Or, optionally, the entry may be multiline. Prompt before
overwriting existing password unless forced.`,
		Run:                   insertPassword,
		Aliases:               []string{"add"},
		DisableFlagsInUseLine: true,
	}
	Echo        bool
	MultiLine   bool
	ForceInsert bool
)

func init() {
	rootCmd.AddCommand(insertCmd)

	insertCmd.Flags().BoolVarP(&Echo, "echo", "e", false, "Echo password back to console")
	insertCmd.Flags().BoolVarP(&MultiLine, "multiline", "m", false, "Multiline input")
	insertCmd.Flags().BoolVarP(&ForceInsert, "force", "f", false, "Overwrite existing password without prompt")
}

func insertPassword(cmd *cobra.Command, args []string) {
	root := util.GetPasswordStore() + "/" + args[0] + ".gpg"
	if f, e := os.Stat(root); !os.IsNotExist(e) && !f.IsDir() {
		if ForceInsert || util.YesNo(fmt.Sprintf("An entry already exists for %s. Overwrite it?", args[0])) {
			if err := os.Remove(root); err != nil {
				fmt.Println(err)
			}
		} else {
			os.Exit(1)
		}
	}
	if MultiLine {
		fmt.Printf("Enter contents of %s and press Ctrl+D when finished:\n\n", args[0])
		encryptMultiLine(root)
	} else {
		encryptPassword(enterPassword(args[0])+"\n", root)
	}
}

func getRecepientOptsArray() []string {
	opts := []string{}
	idFile := util.GetPasswordStore() + "/.gpg-id"
	i, e := os.Open(idFile)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	s := bufio.NewScanner(i)
	for s.Scan() {
		opts = append(opts, "-r")
		opts = append(opts, s.Text())
	}
	return opts
}

func readPassword() (string, error) {
	t := terminal.NewTerminal(os.Stdin, "")
	oldState, _ := terminal.MakeRaw(0)
	defer terminal.Restore(0, oldState)
	if Echo {
		p, e := t.ReadLine()
		return p, e
	}
	p, e := t.ReadPassword("")
	return p, e
}

func enterPassword(name string) string {
	fmt.Printf("Enter password for %s: ", name)
	pass, e := readPassword()
	if e != nil {
		fmt.Println(e)
	}
	fmt.Printf("Retype password for %s: ", name)
	passAgain, e := readPassword()
	if e != nil {
		fmt.Println(e)
	}
	if pass != passAgain {
		fmt.Println("Error: the entered passwords do not match.")
		os.Exit(1)
	}
	return pass
}

func encryptPassword(pass, file string) {
	gpg := exec.Command("gpg",
		"-e", "-o", strings.ReplaceAll(file, " ", `\ `),
		"--quiet", "--yes", "--compress-algo=none", "--no-encrypt-to")
	for _, r := range getRecepientOptsArray() {
		gpg.Args = append(gpg.Args, r)
	}
	stdin, err := gpg.StdinPipe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, pass)
	}()
	os.MkdirAll(filepath.Dir(file), 0755)
	if err := gpg.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func encryptMultiLine(file string) {
	scanner := bufio.NewScanner(os.Stdin)
	pass := []string{}
	for scanner.Scan() {
		pass = append(pass, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	encryptPassword(strings.Join(pass, "\n"), file)
}
