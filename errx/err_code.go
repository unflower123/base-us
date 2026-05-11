package errx

/**
 * @Author: guyu
 * @Desc:
 * @Date: 2025/4/18
 */

// SUCCESS
const OK uint32 = 200

// GLOBAL ERROR
const (
	//MERCHANT
	ErrMerchantAppidNotExist uint32 = 100100 + iota
	ErrMerchantIdNotExist
	ErrMerchantBankListIsEmpty
	ErrMerchantBalanceNotEnough
	ErrMerchantWorkloadOverLimited
	ErrMerchantDailyMaxTxCount
	ErrMerchantTxMinAmount
	ErrMerchantTxMaxAmount
	ErrMerchantAppsecretEmpty
	ErrMerchantStatusDisabled
	ErrMerchantCurrencyIllgal
	ErrMerchantIpIllgal
	ErrMerchantConfigEmpty
	ErrMerchantBalanceEmpty
	ErrMerchantChangeBalanceTypeErr
	ErrMerchantChangeBalancePayTypeErr
	ErrMerchantChangeBalanceFieldErr
	ErrMerchantBankWeightIsZero
	//PAYOUT ORDER
	ErrPayoutOrderExist
	ErrPayoutOrderNotExist
	ErrPayoutOrderStatusErr
	ErrPayoutOrderStatusCallbackErr
	ErrPayoutOrderAmountIsZero
	//MERCHANT Balance
	ErrMerchantBalanceIllegal
)
