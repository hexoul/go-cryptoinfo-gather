package main

import (
	"os"
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
