package main

import (
	"testing"
	"time"

	kucoin "github.com/eeonevision/kucoin-go"
	abcc "github.com/hexoul/go-abcc"
	coinsuper "github.com/hexoul/go-coinsuper"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var (
	testAccessKey = map[string]string{
		"kucoin":    "YOUR_ACCESS_KEY",
		"coinsuper": "YOUR_ACCESS_KEY",
		"abcc":      "YOUR_ACCESS_KEY",
	}
	testSecretKey = map[string]string{
		"kucoin":    "YOUR_SECRET_KEY",
		"coinsuper": "YOUR_SECRET_KEY",
		"abcc":      "YOUR_SECRET_KEY",
	}

	kucoinClient    *kucoin.Kucoin
	abccClient      *abcc.Client
	coinsuperClient *coinsuper.Client
)

func init() {
	kucoinClient = kucoin.New(testAccessKey["kucoin"], testSecretKey["kucoin"])
	abccClient = abcc.GetInstanceWithKey(testAccessKey["abcc"], testSecretKey["abcc"])
	coinsuperClient = coinsuper.GetInstanceWithKey(testAccessKey["coinsuper"], testSecretKey["coinsuper"])
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
	if bal, err := kucoinClient.GetCoinBalance("BTC"); err != nil {
		t.FailNow()
	} else {
		t.Logf("%f %f\n", bal.Balance, bal.FreezeBalance)
	}
}

func TestKucoinListMergedDealtOrders(t *testing.T) {
	if _, err := kucoinClient.ListMergedDealtOrders("ETH-BTC", "BUY", 20, 1, 0, 0); err != nil {
		t.FailNow()
	}
}
