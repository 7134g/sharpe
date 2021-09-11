package conver

import (
	"fmt"
	"strconv"
)

func ReservedFour(f float64) float64 {
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", f), 64)
	return value
}
