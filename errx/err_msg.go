package errx

import "fmt"

/**
 * @Author: guyu
 * @Desc:
 * @Date: 2025/4/18
 */
var message map[uint32]string

func MapErrMsg(errcode uint32) string {
	if msg, ok := message[errcode]; ok {
		return msg
	} else {
		return fmt.Sprintf("errcode:%d", errcode)
	}
}

func init() {
	message = make(map[uint32]string)
	message[OK] = "SUCCESS"
	//Merchant Error
	message[ErrMerchantAppidNotExist] = "Merchant appid is not exist"
	message[ErrMerchantIdNotExist] = "Merchant Id is not exist"
	message[ErrMerchantBankListIsEmpty] = "Merchant bank list is empty"
	message[ErrMerchantBalanceNotEnough] = "Merchant balance is not enough"
	message[ErrMerchantWorkloadOverLimited] = "Merchant workload is over limited"
	message[ErrMerchantDailyMaxTxCount] = "Merchant daily max tx count over"
	message[ErrMerchantTxMinAmount] = "Merchant tx min amount is too less"
	message[ErrMerchantTxMaxAmount] = "Merchant tx max amount is too large"
	message[ErrMerchantAppsecretEmpty] = "Merchant appsecret is empty"
	message[ErrMerchantStatusDisabled] = "Merchant status disabled"
	message[ErrMerchantCurrencyIllgal] = "Merchant currency illgal"
	message[ErrMerchantIpIllgal] = "Merchant ip illgal"
	message[ErrMerchantConfigEmpty] = "Merchant paymethod config  is empty"
	message[ErrMerchantBalanceEmpty] = "Merchant balance is empty"
	message[ErrMerchantChangeBalanceTypeErr] = "Merchant change balance type error"
	message[ErrMerchantChangeBalancePayTypeErr] = "Merchant change balance paytype error"
	message[ErrMerchantChangeBalanceFieldErr] = "Merchant change balance field error"
	message[ErrMerchantBankWeightIsZero] = "Merchant bank weight is zero"
	//PayoutOrder Error
	message[ErrPayoutOrderExist] = " payout orderno exist"
	message[ErrPayoutOrderNotExist] = "payout orderno not exist"
	message[ErrPayoutOrderStatusErr] = "payout order status error"
	message[ErrPayoutOrderStatusCallbackErr] = "payout order callback status error"
	message[ErrPayoutOrderAmountIsZero] = "payout order tx amount is  zero"
	//Merchant Balance
	message[ErrMerchantBalanceIllegal] = "merchant balance change amount illegal"
}
