package authx

import "testing"

func TestGetSecret(t *testing.T) {
	secret := GetSecret()
	if len(secret) != 16 {
		t.Errorf("Expected secret length 16, got %d", len(secret))
	}

	// 检查 secret 是否只包含大写字母
	for _, char := range secret {
		if char < 'A' || char > 'Z' {
			t.Errorf("Secret contains invalid character: %c", char)
		}
	}
}

func TestVerifyCode(t *testing.T) {
	secret := "JBSWY3DPEHPK3PXP" // 这是一个测试用的有效 secret

	// 获取当前时间对应的 code
	code, err := getCode(secret, 0)
	if err != nil {
		t.Fatalf("Failed to get code: %v", err)
	}

	// 验证 code 是否正确
	ok, err := VerifyCode(secret, code)
	if err != nil {
		t.Fatalf("Failed to verify code: %v", err)
	}
	if !ok {
		t.Error("Expected code verification to succeed, but it failed")
	}

	// 测试无效 code
	invalidCode := int32(123456)
	ok, err = VerifyCode(secret, invalidCode)
	if err == nil || ok {
		t.Error("Expected code verification to fail with invalid code, but it succeeded")
	}
}
