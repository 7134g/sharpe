package execute

import (
	"fmt"
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
	start := time.Now()
	ch := make(chan struct{}, 1)
	go func() {
		for {
			end := time.Now()
			select {
			case <-ch:
			default:
				runningTime := end.Sub(start)
				bCount, _ := model.FundBaseLen()
				dCount, _ := model.FundDailyLen()
				log.Println(fmt.Sprintf(
					"\n error: %d poolCount: %d baseTableCount: %d fundTableCount: %d\n"+
						"taskInfo: %d taskDaily: %d taskCore: %d time: %s",
					getTaskErrorCount(),
					fundPool.Cap(),
					bCount,
					dCount,
					infoSpider,
					dailySpider,
					coreSpider,
					runningTime.String()))
				time.Sleep(time.Second * 10)
			}
		}
	}()
	Execute()
	ch <- struct{}{}
	close(ch)
}
