/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/8/6 10:22
 */
package bcryptx

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

func GenerateRSAKeyPKCS(bits int) (privateKeyEncodedStr, publicKeyEncodedStr string, err error) {
	fmt.Printf("Generating RSA key pair with %d bits...\n", bits)

	// 1. 生成 RSA 私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate RSA private key: %w", err)
	}

	// 2. 将私钥编码为 PEM 格式 (PKCS#8)
	// x509.MarshalPKCS8PrivateKey 是现代Go推荐的私钥编码方式
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal private key: %w", err)
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY", // PKCS#8 格式的私钥类型通常是 "PRIVATE KEY"
		Bytes: privateKeyBytes,
	})
	// 对 PEM 格式的私钥字符串进行 Base64 编码
	privateKeyEncodedStr = base64.StdEncoding.EncodeToString(privateKeyPEM)

	// 3. 从私钥中提取公钥，并编码为 PEM 格式 (PKCS#8/SPKI)
	// x509.MarshalPKIXPublicKey 编码公钥为 PKIX (SubjectPublicKeyInfo) 格式，这是 PEM "PUBLIC KEY" 的标准
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal public key: %w", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY", // PKCS#8 格式的公钥类型通常是 "PUBLIC KEY"
		Bytes: publicKeyBytes,
	})
	// 对 PEM 格式的公钥字符串进行 Base64 编码
	publicKeyEncodedStr = base64.StdEncoding.EncodeToString(publicKeyPEM)

	return privateKeyEncodedStr, publicKeyEncodedStr, nil
}

func RSAPrivateKeyFromBase64(base64Key string) (string, error) {

	decodedPrivateKeyBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return "", err
	}

	return string(decodedPrivateKeyBytes), nil
}

func RSAPublicKeyFromBase64(base64Key string) (string, error) {
	decodedPublicKeyBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return "", err
	}

	return string(decodedPublicKeyBytes), nil
}

// RSASignData 使用私钥对数据进行签名
func RSASignData(data []byte, privateKeyBase64 string) (string, error) {

	privateKeyPEM, err := RSAPrivateKeyFromBase64(privateKeyBase64)

	if err != nil {
		return "", fmt.Errorf("failed to decode private key base64 string: %w", err)
	}

	// 解析 PKCS#8 格式的私钥
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil || block.Type != "PRIVATE KEY" {
		return "", fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse PKCS#8 private key: %v", err)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("private key is not an RSA key")
	}

	// 计算数据的 SHA256 哈希
	hasher := sha256.New()
	hasher.Write(data)
	hashed := hasher.Sum(nil)

	// 使用 RSA-SHA256 签名
	signature, err := rsaPrivateKey.Sign(rand.Reader, hashed, crypto.SHA256)
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %v", err)
	}

	// 返回 Base64 编码的签名
	return base64.StdEncoding.EncodeToString(signature), nil
}

// RSAVerifySignature 验证签名
func RSAVerifySignature(data []byte, signatureB64, publicKeyBase64 string) (bool, error) {

	publicKeyPEM, err := RSAPublicKeyFromBase64(publicKeyBase64)

	if err != nil {
		return false, fmt.Errorf("failed to decode public key base64 string : %v", err)
	}

	// 解析 PKCS#8 格式的公钥
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil || block.Type != "PUBLIC KEY" {
		return false, fmt.Errorf("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse PKCS#8 public key: %v", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return false, fmt.Errorf("public key is not an RSA key: got type %T", publicKey)
	}

	// 解码 Base64 签名
	signature, err := base64.StdEncoding.DecodeString(signatureB64)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %v", err)
	}

	// 计算数据的 SHA256 哈希
	hasher := sha256.New()
	hasher.Write(data)
	hashed := hasher.Sum(nil)

	// 验证签名
	err = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hashed, signature)
	if err != nil {
		return false, fmt.Errorf("signature verification failed: %v", err)
	}

	return true, nil
}
