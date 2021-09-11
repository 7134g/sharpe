package execute

import (
	"fmt"
	"net/url"
	"strconv"
)

// CreateInfoURL 生成 info url
func CreateInfoURL(fundType, timeType string) (string, error) {
	u, err := url.Parse(infoURL)
	if err != nil {
		return "", err
	}

	query := u.Query()
	query.Set("ft", fundType)
	query.Set("sc", timeType+"zf")
	query.Set("pn", strconv.Itoa(HttpFundCount))

	u.RawQuery = query.Encode()

	return u.String(), nil
}

// CreateDailyURL 生成当日涨幅 daily url
func CreateDailyURL(code string, page int) (string, error) {
	u, err := url.Parse(dailyURL)
	if err != nil {
		return "", err
	}

	query := u.Query()
	query.Set("code", code)
	query.Set("page", strconv.Itoa(page))
	u.RawQuery = query.Encode()

	return u.String(), nil
}

// CreateCoreUrl 生成夏普率等数据的 url
func CreateCoreUrl(code string) string {
	u := fmt.Sprintf(coreURL, code)
	return u
}
