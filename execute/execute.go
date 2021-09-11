package execute

import (
	"errors"
	"log"
	"sharpe/model"
	"sharpe/pool"
	"time"
)

var (
	fundPool *pool.Pool
)

func init() {
	fundPool = pool.NewPool(POOLSIDE, false, time.Second*5)
}

func Execute() {
	addTaskFundBase()
	fundPool.Wait()
}

func runBase(_ []interface{}) {
	for _, fundType := range fundTypes {
		for _, timeType := range timeTypes {
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
				// 获取每日数据
				go addTaskFundDaily(fb.Code, fb.Name, 0)

				// 获取核心数据
				fc, err := runFundCore(fb.Code)
				if err != nil {
					log.Println(err)
					continue
				}
				// 存储基本信息
				fb.FundType = fundType
				fb.FundCore = *fc
				if err := fb.Create(); err != nil {
					log.Println(err)
					continue
				}
			}

		}
	}

}

func runFundDaily(params []interface{}) {
	code := params[0].(string)
	name := params[1].(string)
	page := params[2].(int)

	u, err := CreateDailyURL(code, page)
	if err != nil {
		log.Println(err)
		return
	}

	fdSheet, err := GetFundDaily(u)
	if errors.Is(err, NoDataError) {
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

	//for _, fd := range fds {
	//	if err := fd.Upsert(); err != nil {
	//		log.Println(err)
	//		return
	//	}
	//}

	if err := model.TransactionFundDaily(fds); err != nil {
		log.Println(err)
		return
	}

	// 下一页
	go addTaskFundDaily(code, name, page)

	return
}

func runFundCore(code string) (*model.FundCore, error) {
	u := CreateCoreUrl(code)
	fcSheet, err := GetFundCoreData(u)
	if err != nil {
		return nil, err
	}

	fc := model.LoadFundCore(fcSheet)

	return fc, nil
}
