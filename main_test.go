package main

import (
	"testing"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func TestGit(t *testing.T) {
	r, err := git.PlainOpen("./")
	if err != nil {
		t.Fatal(err)
	}
	w, err := r.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	if _, err = w.Commit("test go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "hexoul",
			Email: "crosien@gmail.com",
			When:  time.Now(),
		},
		All: true,
	}); err != nil {
		t.Fatal(err)
	}
	if err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: "hexoul",
			Password: "",
		},
	}); err != nil {
		t.Fatal(err)
	}
}
