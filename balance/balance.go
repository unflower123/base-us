package balance

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"math/big"
)

func GetCost(amount int64, rate int64, fee int64) (cost int64, result string) {

	left := big.NewFloat(float64(amount * rate))
	right := big.NewFloat(float64(10000))
	calc := new(big.Float).Quo(left, right)

	// 直接截断小数部分
	calcFloat64, _ := calc.Float64()
	cost = int64(calcFloat64) + fee

	resultFloat, _ := calc.Float64()
	result = fmt.Sprintf("%.2f", resultFloat+float64(fee))
	return
}

func CalculateRate(numerator, denominator int64) float64 {
	if denominator != 0 {
		rate := (float64(numerator) / float64(denominator)) * 100
		return math.Round(rate*100) / 100
	}
	return 0
}

func GetFee(amount int64, rate int64, fee int64) (cost int64, costStr string) {
	_amount := decimal.NewFromInt(amount)
	_rate := decimal.NewFromInt(rate).Div(decimal.NewFromInt(1e6))
	fixedFee := decimal.NewFromInt(fee)

	result := _amount.Mul(_rate).Add(fixedFee)
	cost = result.BigInt().Int64()
	return cost, fmt.Sprintf("%d", cost)
}
