package convertx

import (
	"fmt"
	"github.com/shopspring/decimal"
	"reflect"
	"strconv"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

func toDecimal[T Number](val T) decimal.Decimal {
	switch v := any(val).(type) {
	case int, int8, int16, int32, int64:
		return decimal.NewFromInt(reflect.ValueOf(v).Int())
	case uint, uint8, uint16, uint32, uint64:
		return decimal.NewFromUint64(reflect.ValueOf(v).Uint())
	case float32, float64:
		return decimal.NewFromFloat(reflect.ValueOf(v).Float())
	default:
		return decimal.Decimal{}
	}
}

func WithHundredCarry[T Number](val T, precision int32) T {
	str := WithHundredCarryToString(val, precision)
	_val, _ := strconv.ParseFloat(str, 64)
	return T(_val)
}

func WithCarryToString[T Number](val T, dVal int64, precision int32) string {
	valStr := fmt.Sprintf("%v", val)
	valDecimal, _ := decimal.NewFromString(valStr)
	return valDecimal.Div(decimal.NewFromInt(dVal)).StringFixed(precision)
}

func WithHundredCarryToString[T Number](val T, precision int32) string {
	valDecimal := toDecimal(val)
	_carry := decimal.NewFromInt(1e8)
	result := valDecimal.Div(_carry)
	return result.Truncate(precision).String()
}

func SubWithHundredCarry[T Number](src, sub T, precision int32) T {
	srcDecimal := toDecimal(src)
	subDecimal := toDecimal(sub)
	_carry := decimal.NewFromInt(1e8)
	result := srcDecimal.Sub(subDecimal).Div(_carry)
	_val, _ := strconv.ParseFloat(result.StringFixed(precision), 64)
	return T(_val)
}

func SubWithHundredCarryToString[T Number](src, sub T, precision int32) string {
	srcDecimal := toDecimal(src)
	subDecimal := toDecimal(sub)
	_carry := decimal.NewFromInt(1e8)
	result := srcDecimal.Sub(subDecimal).Div(_carry)
	return result.Truncate(precision).String()
}

func AddWithHundredCarry[T Number](src, add T, precision int32) T {
	srcDecimal := toDecimal(src)
	addDecimal := toDecimal(add)
	_carry := decimal.NewFromInt(1e8)
	result := srcDecimal.Add(addDecimal).Div(_carry)
	_val, _ := strconv.ParseFloat(result.StringFixed(precision), 64)
	return T(_val)
}

func AddWithHundredCarryToString[T Number](src, add T, precision int32) string {
	srcDecimal := toDecimal(src)
	addDecimal := toDecimal(add)
	_carry := decimal.NewFromInt(1e8)
	result := srcDecimal.Add(addDecimal).Div(_carry)
	return result.StringFixed(precision)
}

func Div[T Number](src, sub T, precision int32) T {
	srcDecimal := toDecimal(src)
	var zero T
	if sub == zero {
		_val, _ := strconv.ParseFloat(srcDecimal.Truncate(precision).String(), 64)
		return T(_val)
	}
	subDecimal := toDecimal(sub)
	result := srcDecimal.Div(subDecimal)
	_val, _ := strconv.ParseFloat(result.Truncate(precision).String(), 64)
	return T(_val)
}

func DivToString[T Number](src, sub T, precision int32) string {
	srcDecimal := toDecimal(src)
	var zero T
	if sub == zero {
		return srcDecimal.StringFixed(precision)
	}
	subDecimal := toDecimal(sub)
	result := srcDecimal.Div(subDecimal)
	return result.StringFixed(precision)
}

func Mul[T Number](src, sub T, precision int32) T {
	srcDecimal := toDecimal(src)
	subDecimal := toDecimal(sub)
	result := srcDecimal.Mul(subDecimal)
	_val, _ := strconv.ParseFloat(result.StringFixed(precision), 64)
	return T(_val)
}

func MulToInt64[T Number](src, sub T, precision int32) int64 {
	srcDecimal := toDecimal(src)
	subDecimal := toDecimal(sub)
	result := srcDecimal.Mul(subDecimal)
	//_val, _ := strconv.ParseFloat(result.StringFixed(precision), 64)
	return result.Round(precision).IntPart()
}

func MulToString[T Number](src, sub T, precision int32) string {
	srcDecimal := toDecimal(src)
	subDecimal := toDecimal(sub)
	result := srcDecimal.Mul(subDecimal)
	return result.Truncate(precision).String()
}

func AddWith[T Number](src, add T, precision int32) T {
	srcDecimal := toDecimal(src)
	addDecimal := toDecimal(add)
	result := srcDecimal.Add(addDecimal)
	_val, _ := strconv.ParseFloat(result.StringFixed(precision), 64)
	return T(_val)
}
