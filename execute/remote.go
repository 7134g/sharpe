package execute

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"regexp"
	"sharpe/model"
	"strings"
)

var (
	headerInfo = map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36",
		"Host":       "fund.eastmoney.com",
		"Referer":    "http://fund.eastmoney.com/data/fundranking.html",
		"Accept":     "*/*",
	}
	headerDaily = map[string]string{
		"User-Agent": "Mozilla/5.0(Windows NT 10.0; Win64; x64)AppleWebKit/537.36(KHTML,like Gecko)Chrome/70.0.3538.77Safari/537.36",
		"Host":       "fund.eastmoney.com",
	}
	headerCore = map[string]string{
		"User-Agent": "Mozilla/5.0(Windows NT 10.0; Win64; x64)AppleWebKit/537.36(KHTML,like Gecko)Chrome/70.0.3538.77Safari/537.36",
		"Host":       "fundf10.eastmoney.com",
	}
)

var (
	NoDataError      = errors.New("no data")
	HttpTooFastError = errors.New("too fast")
)

// GetFundInfo 获取基金基础信息
func GetFundInfo(u string) (*model.FundInfoSheet, error) {
	bs, err := get(u, headerInfo)
	if err != nil {
		return nil, err
	}
	cleanBody := bs[15 : len(bs)-1]
	re, err := regexp.Compile(`(\w+):`)
	if err != nil {
		return nil, err
	}
	jsonBody := re.ReplaceAllStringFunc(string(cleanBody), func(s string) string {
		s = fmt.Sprintf(`"%s":`, s[:len(s)-1])
		return s
	})

	fi := &model.FundInfoSheet{}
	if err := json.Unmarshal([]byte(jsonBody), &fi); err != nil {
		return nil, err
	}

	return fi, nil
}

// GetFundDaily 该基金每日数据
func GetFundDaily(u string) (*model.FundDailySheet, error) {
	bs, err := get(u, headerDaily)
	if err != nil {
		return nil, err
	}
	if len(bs) == 0 {
		return nil, HttpTooFastError
	}
	cleanBody := string(bs[12 : len(bs)-1])
	if strings.Contains(cleanBody, "暂无数据") {
		return nil, NoDataError
	}
	re, err := regexp.Compile(`(\w+):`)
	if err != nil {
		return nil, err
	}
	jsonBody := re.ReplaceAllStringFunc(cleanBody, func(s string) string {
		s = fmt.Sprintf(`"%s":`, s[:len(s)-1])
		return s
	})

	fd := &model.FundDailySheet{}
	if err := json.Unmarshal([]byte(jsonBody), &fd); err != nil {
		return nil, err
	}

	return fd, nil
}

// GetFundCoreData 获取该基金核心数据
func GetFundCoreData(u string) (*model.FundCoreSheet, error) {
	bs, err := get(u, headerCore)
	if err != nil {
		return nil, err
	}

	if len(bs) == 0 {
		return nil, HttpTooFastError
	}

	// 解析html
	doc, err := htmlquery.Parse(bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}

	nodes := htmlquery.Find(doc, `//table[@class="fxtb"]//tr`)
	if nodes == nil {
		return nil, errors.New("nodes is nil")
	}

	fc := &model.FundCoreSheet{
		SharpeRatios:      make([]string, 0),
		StandardDeviation: make([]string, 0),
	}
	var sheet = make([]*html.Node, 0)
	for i, node := range nodes {
		sheet = htmlquery.Find(node, `./td`)
		if sheet == nil {
			continue
		}
		switch i {
		case 1:
			for i := 1; i < len(sheet); i++ {
				d := htmlquery.InnerText(sheet[i])
				fc.StandardDeviation = append(fc.StandardDeviation, d)
			}
		case 2:
			for i := 1; i < len(sheet); i++ {
				d := htmlquery.InnerText(sheet[i])
				fc.SharpeRatios = append(fc.SharpeRatios, d)
			}
		}

	}

	return fc, nil
}
