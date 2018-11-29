package main

import (
	"os"
	"strconv"
	"strings"
	"time"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var (
	gitName  string
	gitEmail string
	gitID    string
	gitPW    string
)

func init() {
	for _, val := range os.Args {
		arg := strings.Split(val, "=")
		if len(arg) < 2 {
			continue
		} else if arg[0] == "-gitName" {
			gitName = arg[1]
		} else if arg[0] == "-gitEmail" {
			gitEmail = arg[1]
		} else if arg[0] == "-gitID" {
			gitID = arg[1]
		} else if arg[0] == "-gitPW" {
			gitPW = arg[1]
		}
	}
}

func sumStrFloat(s1, s2 string) (sum float64) {
	f1, err1 := strconv.ParseFloat(s1, 64)
	f2, err2 := strconv.ParseFloat(s2, 64)
	if err1 == nil && err2 == nil {
		sum = f1 + f2
	}
	return
}

func toDateStr(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(time.RFC3339)
}

// gitPushChanges commits log changes and pushs it
func gitPushChanges() error {
	if gitID == "" || gitPW == "" {
		return nil
	}

	// Open
	r, err := git.PlainOpen("./")
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	// Commit
	if _, err = w.Commit("log files changed", &git.CommitOptions{
		Author: &object.Signature{
			Name:  gitName,
			Email: gitEmail,
			When:  time.Now(),
		},
		All: true,
	}); err != nil {
		return err
	}

	// Push
	r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: gitID,
			Password: gitPW,
		},
	})
	return nil
}
