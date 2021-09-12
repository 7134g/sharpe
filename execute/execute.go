package execute

import (
	"errors"
	"log"
	"sharpe/model"
	"sharpe/pool"
	"sync/atomic"
	"time"
)

var (
	infoSpider  int64
	dailySpider int64
	coreSpider  int64
	fundPool    *pool.Pool
)

func init() {
	fundPool = pool.NewPool(POOLSIDE, false, time.Second*5)
}

func Execute() {
	addTaskFundInfo()
	fundPool.Wait()
}

func runInfo(_ []interface{}) {
	for _, fundType := range fundTypes {
		for _, timeType := range timeTypes {
			atomic.AddInt64(&infoSpider, 1)
			u, err := CreateInfoURL(fundType, timeType)
			if err != nil {
				log.Println(err)
				continue
			}

			fi, err := GetFundInfo(u)
			if err != nil {
				log.Println(err)
				continue
			}

			fbs, err := model.LoadFundBase(fi)
			if err != nil {
				log.Println(err)
				continue
			}

			// save
			for _, fb := range fbs {
				fb.FundType = fundType

				// 获取每日数据
				//go addTaskFundDaily(fb.Code, fb.Name, 0)

				// 获取核心数据
				go addTaskFundCore(fb)
			}
			atomic.AddInt64(&infoSpider, -1)
		}
	}

}

func runFundDaily(params []interface{}) {
	code := params[0].(string)
	name := params[1].(string)
	page := params[2].(int)

	if page == 1 {
		atomic.AddInt64(&dailySpider, 1)
	}

	u, err := CreateDailyURL(code, page)
	if err != nil {
		log.Println(err)
		return
	}

	fdSheet, err := GetFundDaily(u)
	if errors.Is(err, NoDataError) {
		atomic.AddInt64(&dailySpider, -1)
		log.Printf("任务完成 ===》 代码：%s, 名称：%s, 页数：%d, url: %s",
			code, name, page, u)
		return
	}
	if errors.Is(err, HttpTooFastError) {
		go retryTask(runFundDaily, params)
		return
	}

	if err != nil {
		log.Println(err)
		return
	}

	fdSheet.Name = name
	fdSheet.Code = code
	fds, err := model.LoadFundDaily(fdSheet)
	if err != nil {
		log.Println(err)
		return
	}

	if err := model.TransactionFundDaily(fds); err != nil {
		log.Println(err)
		return
	}

	// 下一页
	go addTaskFundDaily(code, name, page)

	return
}

func runFundCore(params []interface{}) {
	atomic.AddInt64(&coreSpider, 1)
	fb := params[0].(model.FundBase)

	u := CreateCoreUrl(fb.Code)
	fcSheet, err := GetFundCoreData(u)
	if errors.Is(err, HttpTooFastError) {
		go retryTask(runFundCore, params)
		return
	}
	if err != nil {
		log.Println(err)
		return
	}

	fc := model.LoadFundCore(fcSheet)
	fb.FundCore = *fc

	if err := fb.Upsert(); err != nil {
		log.Println(err)
		return
	}
	atomic.AddInt64(&coreSpider, -1)

	return
}
