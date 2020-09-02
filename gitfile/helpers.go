package gitfile

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/mutexkv"
)

func getConfig(m interface{}) GitFileConfig {
	c := m.(*GitFileConfig)
	res := GitFileConfig{
		Branch:  c.Branch,
		Path:    c.Path,
		RepoUrl: c.RepoUrl,
	}
	return res
}

func cloneIfNotExist(c GitFileConfig) error {
	if _, err := os.Stat(c.Path); err != nil {
		if err = clone(c.Path, c.RepoUrl, c.Branch); err != nil {
			return err
		}
	} else if _, err := os.Stat(path.Join(c.Path, ".git")); err != nil {
		os.RemoveAll(c.Path)
		if err = clone(c.Path, c.RepoUrl, c.Branch); err != nil {
			return err
		}
	}

	return nil
}

func clone(dir, repo, branch string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// May already be checked out from another project
	if _, err := os.Stat(fmt.Sprintf("%s/.git", dir)); err != nil {
		if _, err := gitCommand(dir, "clone", "-b", branch, "--", repo, "."); err != nil {
			cloneExistsError := errwrap.Contains(err, "already exists and is not an empty directory")

			if cloneExistsError {
				log.Printf("[WARN] git clone failed because the directory already exists: : %s", err.Error())
				return nil
			}
			return err
		}
	}
	return nil
}

func gitCommand(checkout_dir string, args ...string) ([]byte, error) {
	command := exec.Command("git", args...)
	command.Dir = checkout_dir
	out, err := command.CombinedOutput()
	if err != nil {
		return out, errwrap.Wrapf(fmt.Sprintf("Error while running git %s: {{err}}\nWorking dir: %s\nOutput: %s", strings.Join(args, " "), checkout_dir, string(out)), err)
	} else {
		return out, err
	}
}

func flatten(args ...interface{}) []string {
	ret := make([]string, 0, len(args))

	for _, arg := range args {
		switch arg := arg.(type) {
		default:
			panic("can only handle strings and []strings")
		case string:
			ret = append(ret, arg)
		case []string:
			ret = append(ret, arg...)
		}
	}

	return ret
}

func hashString(v interface{}) int {
	switch v := v.(type) {
	default:
		panic(fmt.Sprintf("unexpectedtype %T", v))
	case string:
		return hashcode.String(v)
	}
}

// This is a global MutexKV for use within this plugin.
var gitfileMutexKV = mutexkv.NewMutexKV()

func lockCheckout(checkout_dir string) {
	gitfileMutexKV.Lock(checkout_dir)
}

func unlockCheckout(checkout_dir string) {
	gitfileMutexKV.Unlock(checkout_dir)
}

func push(checkout_dir string, count int, retry_count, retry_interval int) error {
	if err := pull(checkout_dir); err != nil {
		return errwrap.Wrapf("push error: {{err}}", err)
	}

	if _, err := gitCommand(checkout_dir, "push", "origin", "HEAD"); err != nil {
		if count >= retry_count {
			return errwrap.Wrapf("retry count elapsed: {{err}}", err)
		}

		time.Sleep(time.Duration(retry_interval) * time.Second)
		count++

		return push(checkout_dir, count, retry_count, retry_interval)
	}
	return nil
}

func commit(checkout_dir, commit_message, commit_body string) error {
	if isEmptyString(commit_body) {
		if _, err := gitCommand(checkout_dir, flatten("commit", "-a", "-m", commit_message, "--allow-empty")...); err != nil {
			return err
		}
	} else {
		if _, err := gitCommand(checkout_dir, flatten("commit", "-a", "-m", commit_message, "-m", commit_body, "--allow-empty")...); err != nil {
			return err
		}
	}
	return nil
}

func pull(checkout_dir string) error {
	if _, err := gitCommand(checkout_dir, "pull", "--strategy=ours"); err != nil {
		return err
	}
	return nil
}

func isEmptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
