package conver

import (
	"github.com/araddon/dateparse"
	"strconv"
	"time"
)

func StringPercentToFloat64(s string) float64 {
	s = s[:len(s)-1]
	v, err := strconv.ParseFloat(s, 0)
	if err != nil {
		return 0
	}
	v = v / 100
	return ReservedFour(v)
}

func StringToFloat64(s string) float64 {
	v, err := strconv.ParseFloat(s, 0)
	if err != nil {
		return 0
	}
	return v
}

func StringToTimeObj(s string) time.Time {
	t, err := dateparse.ParseLocal(s)
	if err != nil {
		return time.Now()
	}
	return t
}
