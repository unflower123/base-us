/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/8/7 11:31
 */
package bcryptx

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func MD5SignData(data []byte) string {

	hashInBytes := md5.Sum(data)

	return strings.ToUpper(hex.EncodeToString(hashInBytes[:]))
}

func MD5VerifySignature(data []byte, expectedHash string) bool {

	calculatedHash := MD5SignData(data)

	return calculatedHash == expectedHash
}
