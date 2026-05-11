/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/4/15 14:36
 */
package consts

const (
	ORDER_BLOOM_ID_KEY                  = "order:%s"
	MERCHANT_BLOOM_ID_KEY               = "bloom:merchant:all"
	BLOOM_PAYIN_ORDER_KEY               = "bloom:payin:order"
	BLOOM_PAYOUT_ORDER_KEY              = "bloom:payout:order"
	MERCHANT_ID_HASH_KEY                = "merchant:id:%d"
	MERCHANT_APPID_HASH_KEY             = "merchant:appid:%s"
	MERCHANT_CONFIG_HASH_KEY            = "merchant:config:%d:%d"    // merchant_id, type
	MERCHANT_PAYOUT_TODAY_TX_KEY        = "merchant:today:%d:%s:%d"  // merchant_id, day, type
	MERCHANT_BANK_LIST_SET_KEY          = "merchant:bank:list:%d:%d" // merchant_id, type
	MERCHANT_BALANCE_HASH_KEY           = "merchant:balance:%d"      //merchant_id
	BANK_INFO_HASH_KEY                  = "bank:info:%d"             //bank_id
	MERCHANT_EMAIL_VERIFY_KEY           = "merchant:email:verify_%d:%s"
	MERCHANT_MODIFY_EMAIL_CHECK_KEY     = "merchant:modify_email:check:%d"
	MERCHANT_MODIFY_PASSWORD_CHECK_KEY  = "merchant:modify_password:check:%d"
	MERCHANT_CONFIG_LIST_KEY            = "merchant:config:list:%d:%d"            // merchant_id, type
	MERCHANT_TODAY_TYPE_GROUP_TX_KEY    = "merchant:today:type:group:%d:%s:%d:%s" // merchant_id, day type group_id
	MERCHANT_TYPE_GROUP_WORKLOAD_KEY    = "merchant:type:group:time:%d:%d:%s:%d"  // merchant_id, type group_id time
	BANK_RATE_INFO_BY_TYPE_KEY          = "bank:rate:info:%d:%d"
	BANK_TODAY_TYPE_TX_KEY              = "bank:today:type:%d:%s:%d" // bank_id, day type
	SYS_CONFIG_INFO_KEY                 = "sys:config:info:%s"
	BANK_STATS_PAYIN                    = "bank:stats:payin"
	BANK_STATS_PAYOUT                   = "bank:stats:payout"
	SUMMARY_BANK_WITHDRAW               = "summary:bank:withdraw"
	SUMMARY_MERCHANT_DATA               = "summary:merchant:data"
	SUMMARY_BANK_UPDATED_CURRENCY       = "summary:bank:updated:currency:%s"
	SUMMARY_MERCHANT_UPDATED_CURRENCY   = "summary:merchant:updated:currency:%s"
	MERCHANT_RESET_PASSWORD_CHECK_KEY   = "merchant:reset_password:check:%s"
	MERCHANT_RESET_PASSWORD_LOCK_KEY    = "merchant:reset_password_look_account:%d"
	MERCHANT_RESET_PASSWORD_LOCK_AT_KEY = "merchant:reset_password_look_at:%d"
	PAYOUT_MAUAL_ORDER                  = "payout:manual:order"
	RISK_CONTROL_RULE_KEY               = "risk:control:rule:%d"
	RISK_CONTROL_RULE_DETAIL_KEY        = "risk:control:rule:detail"
)
