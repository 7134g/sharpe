package model

import (
	"sharpe/db"
	"sync"
)

var (
	once  sync.Once
	dbmux *sync.RWMutex
)

func DBInit() {
	once.Do(func() {
		dbmux = &sync.RWMutex{}
		db.SqliteInit()
		_ = db.DBconn.AutoMigrate(&FundBase{})
		_ = db.DBconn.AutoMigrate(&FundDaily{})
	})
}
