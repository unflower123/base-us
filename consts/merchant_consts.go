package consts

import "fmt"

const (
	MERCHANT_USER_TOKEN = "merchant_user_token"
	MERCHANT_USER_INFO  = "merchant_user_info"
)

func GenerateMerchantUserTokenKey(serverName string, merchantUserID uint64) string {
	return fmt.Sprintf("%s:%s:%v", serverName, MERCHANT_USER_TOKEN, merchantUserID)
}

func GenerateMerchantUserInfoKey(serverName string, merchantUserID uint64) string {
	return fmt.Sprintf("%s:%s:%v", serverName, MERCHANT_USER_INFO, merchantUserID)
}
