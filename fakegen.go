package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/winwisely268/sanfood-faker/client"
	"github.com/winwisely268/sanfood-faker/fakegenerators"
	"github.com/winwisely268/sanfood-faker/model"
)

type BulkAccounts struct {
	Accounts []*generated.AccountsInsertInput `json:"accounts" fakesize:"100"`
}

const (
	baseUri = "https://sf.asiatech.dev/v1/graphql"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error cannot load secret")
	}
	token := os.Getenv("HASURA_ADMIN_SECRET")
	authHeader := func(req *http.Request) {
		req.Header.Set("x-hasura-admin-secret", token)
	}
	// Generate 100 accounts
	bulky := fakegenerators.GenAccounts()
	hc := http.DefaultClient
	cl := client.NewClient(hc, baseUri, authHeader)
	mainCtx := context.Background()
	newCtx, cancel := context.WithTimeout(mainCtx, 5*time.Second)
	defer cancel()
	res, err := cl.BulkInsertAccounts(
		newCtx,
		bulky.Accounts,
		authHeader,
	)
	if err != nil {
		log.Fatalf("error bulk insert accounts: %v", err)
	}
	for _, v := range res.InsertAccounts.Returning {
		log.Printf("Created user with id: %s\n", v.UserID)
	}
}
