/*
预计数据量大约460万左右
*/

package model

import (
	"errors"
	"gorm.io/gorm"
	"sharpe/db"
	"time"
)

// 基金每一天数据
type FundDaily struct {
	gorm.Model
	UID                string    `gorm:"uniqueIndex"` // 用Code和NetWorthDate base64值
	Code               string    // 基金编号
	Name               string    `gorm:"index"` // 基金名称
	NetWorthDate       time.Time // 净值日期
	UnitNetWorth       float64   // 单位净值
	CumulativeNetWorth float64   // 累计净值
	DailyGrowthRate    float64   // 日增长率
}

func (d *FundDaily) Show() ([]FundDaily, error) {
	dbmux.RLock()
	defer dbmux.RUnlock()
	daily := make([]FundDaily, 0)
	err := db.DBconn.Where("code = ?", d.Code).Order(
		"net_worth_date desc").Find(&daily).Error
	if len(daily) == 0 {
		return nil, errors.New("can not find code is " + d.Code)
	}
	if err != nil {
		return nil, err
	}
	return daily, err
}

func (d *FundDaily) Create() error {
	dbmux.Lock()
	defer dbmux.Unlock()
	return db.DBconn.Create(d).Error
}

func (d *FundDaily) Delete() error {
	dbmux.Lock()
	defer dbmux.Unlock()
	return db.DBconn.Delete(d).Error
}

func (d *FundDaily) Update() error {
	dbmux.Lock()
	defer dbmux.Unlock()
	return db.DBconn.Model(&FundDaily{}).Where("uid = ?", d.UID).Updates(d).Error
}

func (d *FundDaily) GetByUID() error {
	dbmux.RLock()
	defer dbmux.RUnlock()
	return db.DBconn.Where("uid = ?", d.UID).First(d).Error
}

func (d *FundDaily) Upsert() error {
	err := d.GetByUID()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return d.Create()
	} else {
		return d.Update()
	}
}

// ShowAllFundDaily 展示所有
func ShowAllFundDaily(code string) ([]FundDaily, error) {
	dbmux.RLock()
	defer dbmux.RUnlock()
	fs := make([]FundDaily, 0)
	err := db.DBconn.Where("code = ?", code).Order("net_worth_date DESC").Find(&fs).Error
	if err != nil {
		return nil, err
	}
	return fs, nil
}

// TransactionFundDaily 事务提交
func TransactionFundDaily(fds []FundDaily) error {
	dbmux.Lock()
	defer dbmux.Unlock()
	err := db.DBconn.Transaction(func(tx *gorm.DB) error {
		for _, fd := range fds {
			err := tx.Create(&fd).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
