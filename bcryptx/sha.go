/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/8/6 14:12
 */
package bcryptx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func SHA256SignData(data []byte, secretKeyBase64 string) (signatureEncodedStr string, err error) {

	secretKey, err := base64.StdEncoding.DecodeString(secretKeyBase64)
	if err != nil {
		err = fmt.Errorf("failed to decode secret key: %w", err)
		return
	}
	// 1. 创建一个新的 HMAC-SHA256 哈希器，使用共享密钥
	h := hmac.New(sha256.New, secretKey)

	// 2. 将数据写入哈希器
	h.Write(data)

	// 3. 计算 HMAC 值（消息认证码）
	// Sum(nil) 会返回最终的 HMAC 字节切片
	signature := h.Sum(nil)

	// 4. 对 HMAC 值进行 Base64 编码并返回
	signatureEncodedStr = base64.StdEncoding.EncodeToString(signature)
	return signatureEncodedStr, nil
}

func SHA256VerifySignature(data []byte, signatureEncodedStr, secretKeyBase64 string) (isValid bool) {

	// 1. 解码收到的 Base64 编码的签名
	secretKey, err := base64.StdEncoding.DecodeString(secretKeyBase64)
	if err != nil {
		fmt.Printf("Error decoding received secret key : %v\n", err)
		return false
	}

	// 1. 解码收到的 Base64 编码的签名
	receivedSignature, err := base64.StdEncoding.DecodeString(signatureEncodedStr)
	if err != nil {
		fmt.Printf("Error decoding received signature: %v\n", err)
		return false
	}

	// 2. 使用相同的密钥和数据，重新计算 HMAC 值
	expectedHMAC := hmac.New(sha256.New, secretKey)
	expectedHMAC.Write(data)
	expectedSignature := expectedHMAC.Sum(nil)

	// 3. 比较计算出的 HMAC 值和收到的 HMAC 值
	// hmac.Equal 是一个安全比较函数，可以防止定时攻击
	isValid = hmac.Equal(receivedSignature, expectedSignature)
	return isValid
}
