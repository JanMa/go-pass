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

func initPasswordStore(cmd *cobra.Command, args []string) {
	gpgKeys := getKeys(args)
	root := util.GetPasswordStore()
	if len(Subdir) > 0 {
		root += "/" + strings.TrimRight(Subdir, "/")
	}
	gpgID := root + "/.gpg-id"
	if _, e := os.Stat(gpgID); os.IsExist(e) {
		exitOnError(os.Remove(gpgID))
	}
	f, e := os.Create(gpgID)
	defer f.Close()
	exitOnError(e)
	f.Write([]byte(strings.Join(args, "\n") + "\n"))
	fmt.Printf("Password store initialized for %s\n", strings.Trim(strings.Join(args, ", "), "\n"))
	gitAddFile(gpgID, fmt.Sprintf("Set GPG id to %s.", strings.Trim(strings.Join(args, ", "), "\n")))
	reEncryptDir(root, gpgKeys)
	gitAddFile(root, fmt.Sprintf("Reencrypt password store using new GPG id %s.", strings.Trim(strings.Join(args, ", "), "\n")))
}

func getKeys(recipients []string) []string {
	re := regexp.MustCompile(`sub:[^:]*:[^:]*:[^:]*:([^:]*):[^:]*:[^:]*:[^:]*:[^:]*:[^:]*:[^:]*:[a-zA-Z]*e[a-zA-Z]*:.*`)
	gpg := exec.Command("gpg", "--list-keys", "--with-colons")
	for _, r := range recipients {
		gpg.Args = append(gpg.Args, r)
	}
	k, e := gpg.Output()
	exitOnError(e)
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
		exitOnError(os.Remove(path))
		gpg := exec.Command("gpg2",
			"-e", "-o", path,
			"--quiet", "--yes", "--compress-algo=none", "--no-encrypt-to")
		for _, r := range keys {
			gpg.Args = append(gpg.Args, "-r")
			gpg.Args = append(gpg.Args, r)
		}
		stdin, err := gpg.StdinPipe()
		exitOnError(err)
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, strings.Join(pass, "\n"))
		}()
		os.MkdirAll(filepath.Dir(path), 0755)
		if o, err := gpg.CombinedOutput(); err != nil {
			fmt.Println(string(o))
			os.Exit(1)
		}
	}
}

func reEncryptDir(path string, keys []string) {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if filepath.Ext(path) == ".gpg" {
			reEncryptFile(path, keys)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
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
