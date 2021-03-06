package execute

import (
	"sharpe/model"
	"sharpe/pool"
	"sync"
)

var (
	errorCount = 0
	lock       = &sync.Mutex{}
)

// base表任务
func addTaskFundInfo() {
	task := &pool.Task{}
	task.TaskFunc = runInfo
	task.Param = []interface{}{}
	_ = fundPool.Submit(task)
}

// daily表任务
func addTaskFundDaily(code, name string, page int) {
	page++
	task := &pool.Task{}
	task.TaskFunc = runFundDaily
	task.Param = []interface{}{code, name, page}
	_ = fundPool.Submit(task)
}

func addTaskFundCore(fb model.FundBase) {
	task := &pool.Task{}
	task.TaskFunc = runFundCore
	task.Param = []interface{}{fb}
	_ = fundPool.Submit(task)
}

// 任务重试
func retryTask(f func([]interface{}), params []interface{}) {
	lock.Lock()
	errorCount++
	lock.Unlock()

	task := &pool.Task{}
	task.TaskFunc = f
	task.Param = params
	_ = fundPool.Submit(task)
}

func getTaskErrorCount() int {
	lock.Lock()
	defer lock.Unlock()
	return errorCount
}
