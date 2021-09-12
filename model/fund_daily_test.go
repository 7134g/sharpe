package model

import (
	"log"
	"os"
	"testing"
)

func init() {
	_ = os.Chdir("../")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	DBInit()
}

func TestTransactionFundDaily(t *testing.T) {
	fds, _ := ShowAllFundDaily("161724")
	err := TransactionFundDaily(fds)
	if err != nil {
		t.Fatal(err)
	}
}
