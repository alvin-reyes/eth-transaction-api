package utils

import (
	"fmt"
	"math"
	"strconv"
)

func WeiToEth(weiStr string, decimals int) string {
	wei, _ := strconv.ParseFloat(weiStr, 64)
	ethValue := wei / math.Pow(10, float64(decimals))
	return fmt.Sprintf("%.18f", ethValue)
}
