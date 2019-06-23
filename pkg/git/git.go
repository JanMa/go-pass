package git

import (
	"fmt"
	"os"
	"os/exec"

	"gitlab.com/JanMa/go-pass/pkg/store"
)

func AddFile(path, msg string) (string, error) {
	s, e := store.GetPasswordStore()
	if f, e := os.Stat(s + "/.git"); os.IsNotExist(e) || !f.IsDir() {
		return "", fmt.Errorf("Password-store is not a git repository")
	}
	o, e := exec.Command("git", "-C", s, "add", path).CombinedOutput()
	if e != nil {
		return string(o), e
	}
	status, _ := exec.Command("git", "-C", s, "status", "--porcelain", path).CombinedOutput()
	if len(msg) > 0 && len(status) > 0 {
		sign := "--no-gpg-sign"
		if c, _ := exec.Command("git", "-C", s, "config", "--bool", "--get", "pass.signcommits").CombinedOutput(); string(c) == "true" {
			sign = "-S"
		}
		o, e = exec.Command("git", "-C", s, "commit", sign, "-m", msg).CombinedOutput()
	}

	return string(o), e
}

func RunCommand(args []string) (string, error) {
	s, e := store.GetPasswordStore()
	if f, e := os.Stat(s + "/.git"); os.IsNotExist(e) || !f.IsDir() {
		return "", fmt.Errorf("Password-store is not a git repository")
	}
	git := exec.Command("git", "-C", s)
	for _, a := range args {
		git.Args = append(git.Args, a)
	}
	o, e := git.CombinedOutput()
	return string(o), e
}
