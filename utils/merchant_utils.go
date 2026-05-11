/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/6/24 11:25
 */
package utils

import (
	"base/consts"
	"fmt"
)

func GetAppidRedisKey(appid string) string {
	return fmt.Sprintf(consts.MERCHANT_APPID_HASH_KEY, appid)
}
