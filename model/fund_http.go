package model

import (
	"encoding/base64"
	"errors"
	"github.com/antchfx/htmlquery"
	"sharpe/conver"
	"strings"
)

type FundInfoSheet struct {
	AllNum   int      `json:"allNum"`
	List     []string `json:"datas"`
	FundType string   `json:"fund_type"`
}

// LoadFundBase 加载基金列表基础数据
func LoadFundBase(fi *FundInfoSheet) ([]FundBase, error) {
	incomes := make([]FundBase, 0)
	for _, v := range fi.List {
		fundList := strings.Split(v, `,`)
		if len(fundList) < 16 {
			return nil, errors.New("extractInfo: data not complete")
		}
		d := FundBase{}
		d.FundType = fi.FundType
		for i := 0; i < len(fundList); i++ {
			switch i {
			case 0:
				d.Code = fundList[i]
			case 1:
				d.Name = fundList[i]
			case 2:
			case 3:
				d.UpdatedAt = conver.StringToTimeObj(fundList[i])
			case 4:
				d.UnitNetWorth = conver.ReservedFour(conver.StringToFloat64(fundList[i]))
			case 5:
				d.CumulativeNetWorth = conver.ReservedFour(conver.StringToFloat64(fundList[i]))
			case 6:
				d.DailyGrowthRate = conver.ReservedFour(conver.StringToFloat64(fundList[i]))
			case 7:
				d.Week = conver.ReservedFour(conver.StringToFloat64(fundList[i]))
			case 8:
				d.Month = conver.ReservedFour(conver.StringToFloat64(fundList[i]))
			case 9:
				d.ThreeMonth = conver.ReservedFour(conver.StringToFloat64(fundList[i]))
			case 10:
				d.SixMonth = conver.ReservedFour(conver.StringToFloat64(fundList[i]))
			case 11:
				d.Year = conver.ReservedFour(conver.StringToFloat64(fundList[i]))
			case 12:
				d.TwoYear = conver.ReservedFour(conver.StringToFloat64(fundList[i]))
			case 13:
				d.ThreeYear = conver.ReservedFour(conver.StringToFloat64(fundList[i]))
			case 14:
				d.NowYear = conver.ReservedFour(conver.StringToFloat64(fundList[i]))
			case 15:
				d.AllTime = conver.ReservedFour(conver.StringToFloat64(fundList[i]))
			case 16:
				d.PublishDate = conver.StringToTimeObj(fundList[i])
			}
		}

		incomes = append(incomes, d)
	}
	return incomes, nil
}

type FundDailySheet struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	FundType string `json:"fund_type"`
	Content  string `json:"content"`
}

// LoadFundDaily 加载列表当日数据
func LoadFundDaily(fd *FundDailySheet) ([]FundDaily, error) {
	funds := make([]FundDaily, 0)
	doc, err := htmlquery.Parse(strings.NewReader(fd.Content))
	if err != nil {
		return nil, err
	}
	nodes := htmlquery.Find(doc, `//table[@class="w782 comm lsjz"]//tr`)
	if nodes == nil {
		return nil, errors.New("nodes is nil")
	}
	for _, node := range nodes {
		once := htmlquery.Find(node, `./td/text()`)
		if once == nil {
			continue
		}
		fund := FundDaily{
			Code: fd.Code,
			Name: fd.Name,
		}
		for i := 0; i < len(once); i++ {
			switch i {
			case 0:
				fund.NetWorthDate = conver.StringToTimeObj(htmlquery.InnerText(once[i]))
				uid := fund.Code + fund.NetWorthDate.Format(" 2006-01-02")
				fund.UID = base64.StdEncoding.EncodeToString([]byte(uid))
			case 1:
				fund.UnitNetWorth = conver.StringToFloat64(htmlquery.InnerText(once[i]))
			case 2:
				fund.CumulativeNetWorth = conver.StringToFloat64(htmlquery.InnerText(once[i]))
			case 3:
				dgr := strings.ReplaceAll(htmlquery.InnerText(once[i]), "%", "")
				fund.DailyGrowthRate = conver.StringToFloat64(dgr)
			}
		}
		funds = append(funds, fund)

	}
	return funds, nil
}

type FundCoreSheet struct {
	SharpeRatios      []string
	StandardDeviation []string
}

// LoadFundCore 加载基金核心数据
func LoadFundCore(fc *FundCoreSheet) *FundCore {
	fb := &FundCore{}
	for i, v := range fc.SharpeRatios {
		switch i {
		case 0:
			fb.SharpeRatioYear = conver.StringToFloat64(v)
		case 1:
			fb.SharpeRatioTwoYear = conver.StringToFloat64(v)
		case 2:
			fb.SharpeRatioThreeYear = conver.StringToFloat64(v)
		}
	}

	for i, v := range fc.StandardDeviation {
		switch i {
		case 0:
			fb.StandardDeviationYear = conver.StringPercentToFloat64(v)
		case 1:
			fb.StandardDeviationTwoYear = conver.StringPercentToFloat64(v)
		case 2:
			fb.StandardDeviationThreeYear = conver.StringPercentToFloat64(v)
		}
	}
	return fb
}
