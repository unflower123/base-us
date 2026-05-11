/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/6/5 17:55
 */
package utils

import (
	"fmt"
	"math"
)

func FenToYuan(fen int) string {
	yuan := float64(fen) / 100.0

	return fmt.Sprintf("%.2f", yuan)
}

func FenToYuanFloat64(fen int64) float64 {
	return float64(fen) / 100.0
}

func FenToYuanFloat32(fen int64) float32 {
	return float32(fen) / 100.0
}

func YuanToFen(yuan float32) int64 {
	return int64(yuan * 100)
}
func YuanToFenFloat64(yuan float64) int64 {
	fenFloat := math.Round(yuan * 100)
	return int64(fenFloat)
}
