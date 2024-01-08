package versionx

import (
	"io/ioutil"
	"os"
	"strings"
)

var (
	Version   = "undefined"
	GitHash   = "undefined"
	GitBranch = "undefined"
	BuildDate = "undefined"
)

func init() {
	if _, err := os.Stat("version.txt"); err == nil {
		if dat, err := ioutil.ReadFile("version.txt"); err == nil {
			Version = strings.TrimSpace(string(dat))
		}
	}

	if _, err := os.Stat("commit.txt"); err == nil {
		if dat, err := ioutil.ReadFile("commit.txt"); err == nil {
			GitHash = strings.TrimSpace(string(dat))
		}
	}
	if _, err := os.Stat("buildDate.txt"); err == nil {
		if dat, err := ioutil.ReadFile("buildDate.txt"); err == nil {
			BuildDate = strings.TrimSpace(string(dat))
		}
	}
	if _, err := os.Stat("branch.txt"); err == nil {
		if dat, err := ioutil.ReadFile("branch.txt"); err == nil {
			GitBranch = strings.TrimSpace(string(dat))
		}
	}
}
