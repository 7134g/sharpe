package execute

import "testing"

func TestGetFundInfo(t *testing.T) {
	u := `http://fund.eastmoney.com/data/rankhandler.aspx?op=ph&dt=kf&ft=gp&rs=&gs=0&sc=1yzf&st=desc&pi=1&pn=10000&dx=1`
	fi, err := GetFundInfo(u)
	if err != nil {
		panic(err)
	}
	t.Log(fi)
}

func TestGetFundDaily(t *testing.T) {
	u := `http://fund.eastmoney.com/f10/F10DataApi.aspx?type=lsjz&code=004734&page=1&per=1000`
	fd, err := GetFundDaily(u)
	if err != nil {
		panic(err)
	}
	t.Log(fd)
}

func TestGetFundCoreData(t *testing.T) {
	u := `http://fundf10.eastmoney.com/tsdata_009126.html`
	fd, err := GetFundCoreData(u)
	if err != nil {
		panic(err)
	}
	t.Log(fd)
}
