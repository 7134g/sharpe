package model

import (
	"errors"
	"gorm.io/gorm"
	"sharpe/db"
	"time"
)

// FundBase 基金表
type FundBase struct {
	gorm.Model
	Code               string    `gorm:"uniqueIndex"` // 基金编号
	FundType           string    `gorm:"index"`       // 基金类型
	Name               string    `gorm:"index"`       // 基金名称
	UnitNetWorth       float64   // 单位净值
	CumulativeNetWorth float64   // 累计净值
	DailyGrowthRate    float64   // 日增长率
	Week               float64   // 近一周
	Month              float64   // 近一月
	ThreeMonth         float64   // 近三月
	SixMonth           float64   // 近六月
	Year               float64   // 近一年
	TwoYear            float64   // 近两年
	ThreeYear          float64   // 近三年
	NowYear            float64   // 今年
	AllTime            float64   // 成立以来
	PublishDate        time.Time // 发行时间
	ManagerTime        time.Time // 基金经理管理时间(未实现)

	FundCore FundCore `gorm:"embedded"`
}

type FundCore struct {
	SharpeRatioYear            float64 // 1年夏普率
	SharpeRatioTwoYear         float64 // 2年夏普率
	SharpeRatioThreeYear       float64 // 3年夏普率
	StandardDeviationYear      float64 // 1年标准差
	StandardDeviationTwoYear   float64 // 2年标准差
	StandardDeviationThreeYear float64 // 3年标准差
}

func (b *FundBase) Get() (*FundBase, error) {
	dbmux.RLock()
	defer dbmux.RUnlock()
	fund := &FundBase{}
	err := db.DBconn.Where("code = ?", b.Code).First(fund).Error
	if err != nil {
		return nil, err
	}
	return fund, nil
}

func (b *FundBase) Create() error {
	dbmux.Lock()
	defer dbmux.Unlock()
	return db.DBconn.Create(b).Error
}

func (b *FundBase) Delete() error {
	dbmux.Lock()
	defer dbmux.Unlock()
	return db.DBconn.Delete(b).Error
}

func (b *FundBase) Update() error {
	dbmux.Lock()
	defer dbmux.Unlock()
	return db.DBconn.Model(&FundBase{}).Where("code = ?", b.Code).Updates(b).Error
}

func (b *FundBase) Upsert() error {
	_, err := b.Get()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return b.Create()
	} else {
		return b.Update()
	}
}

func FundBaseLen() (int64, error) {
	dbmux.RLock()
	defer dbmux.RUnlock()
	var count int64
	if err := db.DBconn.Model(&FundBase{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func ShowAllFundBase(page, size int) ([]FundBase, error) {
	dbmux.RLock()
	defer dbmux.RUnlock()
	var err error
	fs := make([]FundBase, 0)
	if page*size == 0 {
		err = db.DBconn.Find(&fs).Error
	} else {
		err = db.DBconn.Limit(size).Offset(page * size).Find(&fs).Error
	}
	if err != nil {
		return nil, err
	}
	return fs, nil
}

func ShowMaxSharpe(key, order string, size int) ([]FundBase, error) {
	dbmux.RLock()
	defer dbmux.RUnlock()

	var tx *gorm.DB
	fs := make([]FundBase, 0)
	if key != "" {
		v := "%" + key + "%"
		tx = db.DBconn.Where("name LIKE ?", v)
	}

	if size != 0 {
		tx = tx.Limit(size)
	}

	order = order + " desc"
	err := tx.Order(order).Find(&fs).Error
	if err != nil {
		return nil, err
	}
	return fs, nil
}
