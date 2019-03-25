package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/JanMa/go-pass/util"
)

// initCmd represents the init command
var (
	initCmd = &cobra.Command{
		Use:   "init",
		Args:  cobra.MinimumNArgs(1),
		Short: "Initialize new password storage and use gpg-id for encryption.",
		Long: `Initialize new password storage and use gpg-id for encryption.
Selectively reencrypt existing passwords using new gpg-id.`,
		Run: initPasswordStore,
	}

	Subdir string
)

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&Subdir, "path", "p", "", "A specific gpg-id or set of gpg-ids is assigned for that specific sub folder of the password store")
}

func initPasswordStore(cmd *cobra.Command, args []string) {
	gpgKeys := getKeys(args)
	root := util.GetPasswordStore()
	if len(Subdir) > 0 {
		root += "/" + Subdir
	}
	gpgID := root + "/.gpg-id"
	if _, e := os.Stat(gpgID); os.IsExist(e) {
		if e := os.Remove(gpgID); e != nil {
			fmt.Println(e)
			os.Exit(1)
		}
	}
	f, e := os.Create(gpgID)
	defer f.Close()
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	f.Write([]byte(strings.Join(args, "\n") + "\n"))
	fmt.Printf("Password store initialized for %s\n", strings.Trim(strings.Join(args, ", "), "\n"))
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if filepath.Ext(path) == ".gpg" {
			reEncryptFile(path, gpgKeys)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getKeys(recipients []string) []string {
	re := regexp.MustCompile(`sub:[^:]*:[^:]*:[^:]*:([^:]*):[^:]*:[^:]*:[^:]*:[^:]*:[^:]*:[^:]*:[a-zA-Z]*e[a-zA-Z]*:.*`)
	gpg := exec.Command("gpg", "--list-keys", "--with-colons")
	for _, r := range recipients {
		gpg.Args = append(gpg.Args, r)
	}
	k, e := gpg.Output()
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	match := re.FindAllSubmatch(k, -1)

	gpgKeys := []string{}
	for _, m := range match {
		gpgKeys = append(gpgKeys, string(m[1]))
	}

	return gpgKeys
}

func reEncryptFile(path string, keys []string) {
	currentKeys, _ := exec.Command("gpg", "-v", "-d", "--list-only", "--keyid-format", "long", path).CombinedOutput()
	if !matchKeys(string(currentKeys), keys) {
		fmt.Printf("%s: reencrypting to %s\n", filepath.Base(path), strings.Join(keys, ", "))
		pass := util.RunCommand("gpg", "-dq", path)
		if e := os.Remove(path); e != nil {
			fmt.Println(e)
		}
		gpg := exec.Command("gpg2",
			"-e", "-o", path,
			"--quiet", "--yes", "--compress-algo=none", "--no-encrypt-to")
		for _, r := range keys {
			gpg.Args = append(gpg.Args, "-r")
			gpg.Args = append(gpg.Args, r)
		}
		stdin, err := gpg.StdinPipe()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, strings.Join(pass, "\n"))
		}()
		os.MkdirAll(filepath.Dir(path), 0755)
		if err := gpg.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func matchKeys(input string, keys []string) bool {
	for _, k := range keys {
		if ret, _ := regexp.MatchString(k, input); !ret {
			return false
		}
	}
	return true
}
