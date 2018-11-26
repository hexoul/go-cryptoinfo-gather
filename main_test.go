package main

import (
	"testing"
	"time"

	kucoin "github.com/eeonevision/kucoin-go"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var (
	testAccessKey = map[string]string{
		"coinsuper": "YOUR_ACCESS_KEY",
		"kucoin":    "YOUR_ACCESS_KEY",
		"abcc":      "YOUR_ACCESS_KEY",
	}
	testSecretKey = map[string]string{
		"coinsuper": "YOUR_SECRET_KEY",
		"kucoin":    "YOUR_SECRET_KEY",
		"abcc":      "YOUR_SECRET_KEY",
	}
)

func init() {

}

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

func TestKucoinBalance(t *testing.T) {
	k := kucoin.New(testAccessKey["kucoin"], testSecretKey["kucoin"])
	if _, err := k.GetCoinBalance("BTC"); err != nil {
		t.FailNow()
	}
}

func TestKucoinListMergedDealtOrders(t *testing.T) {
	k := kucoin.New(testAccessKey["kucoin"], testSecretKey["kucoin"])
	if _, err := k.ListMergedDealtOrders("ETH-BTC", "BUY", 20, 1, 0, 0); err != nil {
		t.FailNow()
	}
}
