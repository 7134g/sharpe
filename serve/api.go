package serve

import (
	"errors"
	"net/http"
	"sharpe/execute"
	"sharpe/model"
)

var flag bool

func BaseInfo(ctx *BaseContext) {
	page, err := ctx.QueryInt("page")
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}
	size, err := ctx.QueryInt("size")
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}

	result := struct {
		FundBases []model.FundBase
		Total     int64
	}{}

	count, err := model.FundBaseLen()
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}

	fbs, err := model.ShowAllFundBase(page, size)
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}

	result.FundBases = fbs
	result.Total = count

	ctx.JSON(http.StatusOK, result)
}

func FundDaily(ctx *BaseContext) {
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusOK, "code is nil")
		return
	}

	fd := model.FundDaily{Code: code}
	fds, err := fd.Show()
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}
	ctx.JSON(http.StatusOK, fds)
}

func PullFund(ctx *BaseContext) {
	if !flag {
		flag = true
		go func() {
			execute.Execute()
			flag = false
		}()
		ctx.JSON(http.StatusOK, "ok")
	} else {
		ctx.JSON(http.StatusOK, "executing")
	}
	return
}

func MaxSharpe(ctx *BaseContext) {
	name := ctx.Query("name")
	fType := ctx.Query("fType")
	tType := ctx.Query("tType")
	if tType == "" {
		ctx.JSON(http.StatusOK, errors.New("tType is nil"))
		return
	}
	size := ctx.QueryIntDefault("size")

	var order string
	switch tType {
	case "1y":
		order = "sharpe_ratio_year"
	case "2y":
		order = "sharpe_ratio_two_year"
	case "3y":
		order = "sharpe_ratio_three_year"
	default:
		ctx.JSON(http.StatusOK, errors.New("tType is error: "+tType))
		return
	}

	switch fType {
	case "gp", "hh", "zq", "zs":
		break
	default:
		ctx.JSON(http.StatusOK, errors.New("fType is error: "+tType))
		return
	}

	fbs, err := model.ShowMaxSharpe(fType, name, order, size)
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}

	var result []struct {
		Name   string
		Code   string
		Sharpe float64
	}

	var sharpe float64
	for _, fb := range fbs {
		switch tType {
		case "1y":
			sharpe = fb.FundCore.SharpeRatioYear
		case "2y":
			sharpe = fb.FundCore.SharpeRatioTwoYear
		case "3y":
			sharpe = fb.FundCore.SharpeRatioThreeYear
		}
		result = append(result, struct {
			Name   string
			Code   string
			Sharpe float64
		}{Name: fb.Name, Code: fb.Code, Sharpe: sharpe})
	}

	ctx.JSON(http.StatusOK, result)
}
