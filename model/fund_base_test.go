package model

import "testing"

func TestTransactionFundBase(t *testing.T) {
	fbs, _ := ShowAllFundBase(1, 10)
	err := TransactionFundBase(fbs)
	if err != nil {
		t.Fatal(err)
	}
}
