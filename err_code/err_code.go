// go install golang.org/x/tools/cmd/stringer@latest

//go:generate stringer -type=ErrCode --linecomment
// stringer -type=ErrCode --linecomment

package err_code

type ErrCode int

func (i ErrCode) Error() string {
	return i.String()
}

func (i ErrCode) GetCode() int32 {
	return int32(i)
}

const (
	ErrCodeSuccess ErrCode = 100000 // success
)

const (
	ErrCodeRepeatOrder                         ErrCode = 400001 // Duplicate order number
	ErrCodeNotifyUrlEMpty                      ErrCode = 400002 // The notification address of the merchant or order is empty
	ErrCodeMerchantNotFound                    ErrCode = 400003 // merchant does not exist
	ErrCodeOrderNotFound                       ErrCode = 400004 // order does not exist
	ErrCodeOrderRecordNotFound                 ErrCode = 400005 // order record not found
	ErrCodeResuestParamCheck                   ErrCode = 400006 // request parameter verification failed
	ErrCodeMerchantPayConfigNotFound           ErrCode = 400007 // merchant pay config does not exist
	ErrCodeMerchantAppSecretCheck              ErrCode = 400008 // merchant AppSecret verification failed
	ErrCodeMerchantCurrencyCheck               ErrCode = 400009 // merchant Currency verification failed
	ErrCodeMerchantIpCheck                     ErrCode = 400010 // merchant Ip verification failed
	ErrCodeMerchantAmountCheck                 ErrCode = 400011 // merchant transaction amount lt 0
	ErrCodeMerchantDailyMaxTxCountCheck        ErrCode = 400012 // Total amount has reached daily limit.
	ErrCodeMerchantStatusCheck                 ErrCode = 400013 // merchant status is disabled
	ErrCodeMerchantWorkloadCheck               ErrCode = 400014 // The volume per minute has reached the limit.
	ErrCodeMerchantBalanceCheck                ErrCode = 400015 // merchant balance verification failed, don't have enough balance
	ErrCodeMerchantBankNotFound                ErrCode = 400016 // merchant choose bank fail
	ErrCodeMerchantTradeMinCheck               ErrCode = 400017 // Transaction amount is below the minimum limit. Please re-enter.
	ErrCodeMerchantTradeMaxCheck               ErrCode = 400018 // Transaction amount exceeds the maximum limit. Please re-enter.
	ErrCodeMerchantOrderNoAlreadyExist         ErrCode = 400019 // the order no is duplicated for this merchant
	ErrCodeMerchantNotifyUrlIsNull             ErrCode = 400020 // merchant notify url is null
	ErrCodeProtoMarshal                        ErrCode = 400021 // proto marshal failed
	ErrCodeProtoUnmarshal                      ErrCode = 400022 // proto unmarshal failed
	ErrCodePrePayIDNotExist                    ErrCode = 400023 // order expired
	ErrCodeMerchantNotMatch                    ErrCode = 400024 // merchantAppid not match
	ErrCodeMerchantCurrencyIsNull              ErrCode = 400025 // merchant currency is null
	ErrCodeMerchantIPNotInWhiteList            ErrCode = 400026 // this IP is not on the whitelist
	ErrCodeMerchantSignatureVerificationFailed ErrCode = 400027 // signature verification failed
	ErrCodeMerchantSignValAssertionFailed      ErrCode = 400028 // sign val assertion failed
	ErrCodeTimestampExpired                    ErrCode = 400029 // Timestamp expired, please set a timestamp within 1 minutes
	ErrCodeMerchantEDBaseDecode                ErrCode = 400030 // Merchant public key base decoding failed
	ErrCodeMerchantWithdrawaSwitchDisabled     ErrCode = 400031 // Merchant withdrawal switch disabled
	ErrCodeBankSwitchDisabled                  ErrCode = 400032 // This bank switch disabled
	ErrCodeBankDailyMaxTxCountCheck            ErrCode = 400033 // Check:Total amount has reached daily limit.
	ErrCodeBankWorkloadCheck                   ErrCode = 400034 // Check:The volume per minute has reached the limit.
	ErrCodeBankTradeMinCheck                   ErrCode = 400035 // Check:Transaction amount is below the minimum limit. Please re-enter.
	ErrCodeBankTradeMaxCheck                   ErrCode = 400036 // Check:Transaction amount exceeds the maximum limit. Please re-enter.
	ErrCodeMerchantRateNotEdit                 ErrCode = 400037 // merchant choose group fail
	ErrCodeBankIsDisabled                      ErrCode = 400038 // Check:merchant choose bank fail
	ErrCodeMerchantBankIsDisabled              ErrCode = 400039 // merchant choose bank fail
	ErrCodeMOrderAmountNotSupportDecimal       ErrCode = 400040 // Order Amount does not support decimals

)

const (
	ErrCodeSystem               ErrCode = 500001 // internal dependency error
	ErrCodeRedis                ErrCode = 500002 // internal dependency error
	ErrCodeKafka                ErrCode = 500003 // internal dependency error
	ErrCodeFormatPB             ErrCode = 500004 // internal dependency error
	ErrCodeRpc                  ErrCode = 500005 // internal dependency error
	ErrCodeDataBase             ErrCode = 500006 // internal dependency error
	ErrCodeBase                 ErrCode = 500007 // internal dependency error
	ErrCodeMerchantRPC          ErrCode = 500008 // internal dependency error
	ErrCodeMerchantUnmarshal    ErrCode = 500009 // internal dependency error
	ErrCodeMethodTypeNotAllowed ErrCode = 500010 // method type not allowed
	ErrBankServerIPEmpty        ErrCode = 500011 // internal dependency error
	ErrBankServerRpc            ErrCode = 500012 // internal dependency error
)

const (
	ErrCodeMerchant ErrCode = 20000 // merchant err
)

const (
	ErrCodeAdmin ErrCode = 30000 // admin err
)
