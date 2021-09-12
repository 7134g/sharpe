package execute

var (
	// 基金类型
	fundTypes = []string{"gp", "hh", "zq", "zs"}
	// 时间类型
	timeTypes = []string{"1y", "3y", "6y", "1n", "2n"}
	// 单次获取基金数据量
	HttpFundCount = 10000
	// 代理池大小
	POOLSIDE int32 = 300

	infoURL = "http://fund.eastmoney.com/data/rankhandler.aspx?" +
		"op=ph&dt=kf&ft=gp&rs=&gs=0&sc=1yzf&st=desc&pi=1&pn=10000&dx=1"
	dailyURL = "https://fundf10.eastmoney.com/F10DataApi.aspx?" +
		"type=lsjz&code=004734&page=1&per=1000"
	coreURL = "http://fundf10.eastmoney.com/tsdata_%s.html"
)
