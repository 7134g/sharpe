package execute

import (
	"log"
	"os"
	"sharpe/model"
	"testing"
	"time"
)

func init() {
	_ = os.Chdir("../")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	model.DBInit()
}

func TestExecute(t *testing.T) {
	ch := make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-ch:
			default:
				log.Println("error count", getTaskErrorCount())
				time.Sleep(time.Second)
			}
		}
	}()
	Execute()
	ch <- struct{}{}
	close(ch)
}
