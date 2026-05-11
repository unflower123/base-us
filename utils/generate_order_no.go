/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/5/20 10:15
 */
package utils

import (
	"base/consts"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const characterSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateOrderNo(currency string, merchantId int64, orderType int) string {

	countryMark := consts.CurrencyMappingCountry[currency]
	orderMark := consts.OrderMappingMark[orderType]

	midStr := strconv.Itoa(int(merchantId))

	mid := midStr[len(midStr)-3:]

	uuidMark, _ := GenerateRandomString(10)

	t := time.Now()
	mdTime := t.UTC().Format("0102")

	orderNo := fmt.Sprintf("%s%d%s%v%s", countryMark, orderMark, mid, mdTime, uuidMark)
	//orderNo := fmt.Sprintf("%s%d%d%v%s", countryMark, orderMark, merchantId, mdTime, uuidMark)

	return strings.ToUpper(orderNo)
}

func GenerateRandomString(length int) (string, error) {
	// 创建一个字节切片来存储结果。
	// Create a byte slice to store the result.
	result := make([]byte, length)
	// 获取字符集的长度。
	// Get the length of the character set.
	charsetLength := big.NewInt(int64(len(characterSet)))

	// 循环生成每个字符。
	// Loop to generate each character.
	for i := 0; i < length; i++ {
		// 生成一个安全的随机索引。
		// Generate a secure random index.
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		// 根据随机索引从字符集中选择一个字符。
		// Select a character from the character set based on the random index.
		result[i] = characterSet[randomIndex.Int64()]
	}

	return string(result), nil
}
