/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/9/3 16:26
 */
package utils

func CalculateSuccessRate(successCount, totalCount int64) float64 {
	// 如果总数为0，为避免除以零的错误，成功率被视为0。
	// 我们返回格式化的字符串 "0.00%"。
	if totalCount == 0 {
		return 0
	}
	// 执行计算并转换为百分比。
	rate := float64(successCount) / float64(totalCount) * 100

	return rate
}

// 计算平均值，并返回一个格式化的字符串。
// Parameters:
//
//	value: 计算值（被除数）。
//	divisor: 被除值（除数）。
//
// Returns:
//
//	一个表示平均数的字符串，格式化为两位小数（例如："30.25"）。
func CalculateAverage(value, divisor int64) float64 {
	// 处理被除数为零的边界情况，以避免除以零的错误。
	if divisor == 0 {
		return 0
	}

	average := float64(value) / float64(divisor)
	return average
}
