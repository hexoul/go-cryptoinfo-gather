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

	testClients Clients
)

func init() {
	testClients.kucoin = kucoin.New(testAccessKey["kucoin"], testSecretKey["kucoin"])
	testClients.abcc = abcc.GetInstanceWithKey(testAccessKey["abcc"], testSecretKey["abcc"])
	testClients.coinsuper = coinsuper.GetInstanceWithKey(testAccessKey["coinsuper"], testSecretKey["coinsuper"])
}

func TestKucoinBalance(t *testing.T) {
	if bal, err := testClients.kucoin.GetCoinBalance("USDT"); err != nil {
		t.FailNow()
	} else {
		t.Logf("%f %f\n", bal.Balance, bal.FreezeBalance)
	}
}

func TestKucoinListMergedDealtOrders(t *testing.T) {
	if _, err := testClients.kucoin.ListMergedDealtOrders("ETH-BTC", "BUY", 20, 1, 0, 0); err != nil {
		t.FailNow()
	}
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
