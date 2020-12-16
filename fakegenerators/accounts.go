package fakegenerators

import (
	"github.com/brianvoe/gofakeit/v5"

	"github.com/winwisely268/sanfood-faker/fakehelper"
	generated "github.com/winwisely268/sanfood-faker/model"
	"github.com/winwisely268/sanfood-faker/utilities"
)

type BulkAccounts struct {
	Accounts []*generated.AccountsInsertInput `json:"accounts" fakesize:"100"`
}

func GenAccounts() *BulkAccounts {
	var bulky BulkAccounts
	gofakeit.Seed(utilities.CurrentTimestamp())
	gofakeit.AddFuncLookup(fakehelper.FakeGenRole())
	gofakeit.AddFuncLookup(fakehelper.FakeMailGen())
	gofakeit.Struct(&bulky)
	return &bulky
}

