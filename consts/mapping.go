package consts

func CreateTypeMapping(createType int64) string {
	switch createType {
	case CREATE_TYPE_SYSTEM:
		return "API"
	case CREATE_TYPE_MANUAL:
		return "Manual Entry"
	case CREATE_TYPE_MANUAL_CALLBACK:
		return "Manual Callback"
	default:
		return ""
	}
}

func BankAccountTypeMapping(bankAccountType int64) string {
	switch bankAccountType {
	case BANK_ACCOUNT_TYPE_PAYIN_NATIVE:
		return "Intent"
	case BANK_ACCOUNT_TYPE_PAYIN_QRCODE:
		return "MQR"
	case BANK_ACCOUNT_TYPE_PAYIN_TRANSFER:
		return "Transfer"
	case BANK_ACCOUNT_TYPE_PAYOUT:
		return "Pay out"
	case BANK_ACCOUNT_TYPE_SETTLEMENT:
		return "Settlement"
	case BANK_ACCOUNT_TYPE_CHANGE:
		return "Balance Management"
	default:
		return ""
	}
}

func PayStatusMapping(payStatus int64) string {
	switch payStatus {
	case PAY_STATUS_SUCCESS:
		return "Success"
	case PAY_STATUS_FAILDE:
		return "Failure"
	case PAY_STATUS_ABNORMAL:
		return "Abnormal"
	case PAY_STATUS_NOT:
		return "Pending"
	case PAY_STATUS_WAIT_AUDIT:
		return "Not Review"
	case PAY_STATUS_REJECT:
		return "Reject"
	case PAY_STATUS_REFUND:
		return "Reversed"
	case PAY_STATUS_COMPLAIN:
		return "complaint"
	default:
		return ""
	}
}

func NotifyStatusMapping(notifyStatus int64) string {
	switch notifyStatus {
	case NOTIFY_STATUS_SUCCESS:
		return "Notified"
	case NOTIFY_STATUS_FAILED:
		return "Notify failed"
	case NOTIFY_STATUS_NOT:
		return "Not Notified"
	default:
		return ""
	}
}

func PayoutStatusMapping(status int64) string {
	switch status {
	case PAYOUT_STATUS_SUCCESS:
		return "Success"
	case PAYOUT_STATUS_FAILED:
		return "Failure"
	case PAYOUT_STATUS_PENDING:
		return "Pending"
	case PAYOUT_STATUS_REFUNDED:
		return "Refunded"
	default:
		return ""
	}
}

func BalanceTradeCommandMapping(status int64) string {
	switch status {
	case BALANCE_TRADE_COMMAND_SUCCESS:
		return "Success"
	case BALANCE_TRADE_COMMAND_ADDITIONAL_ORDER:
		return "Record Missing Payin Orders"
	case BALANCE_TRADE_COMMAND_COMPLAINT:
		return "Complaint"
	case BALANCE_TRADE_COMMAND_REVERSAL:
		return "Refunded"
	case BALANCE_TRADE_COMMAND_INCRE:
		return "Increase"
	case BALANCE_TRADE_COMMAND_DECR:
		return "Decrease"
	default:
		return ""
	}
}

func BalanceFreezeMapping(status int64) string {
	switch status {
	case BALANCE_FREEZE_STATE:
		return "block"
	case BALANCE_UNFREEZE_STATE:
		return "unblock"
	default:
		return ""
	}
}

func BalanceFreezeTradeCommandMapping(status int64) string {
	switch status {
	case BALANCE_FREEZE_TRADE_COMMAND_PAYOUT:
		return "Payout Initiated"
	case BALANCE_FREEZE_TRADE_COMMAND_WITHDRAWAL:
		return "Payout Success"
	case BALANCE_FREEZE_TRADE_COMMAND_MERCHANT_WITHDRAWAL_INITIATED:
		return "Merchant Withdraw Initiated"
	case BALANCE_FREEZE_TRADE_COMMAND_MERCHANT_WITHDRAWAL_REJECTED:
		return "Merchant Withdraw Rejected"
	case BALANCE_FREEZE_TRADE_COMMAND_MERCHANT_WITHDRAWAL_SUCCESS:
		return "Merchant Withdraw Success"
	case BALANCE_FREEZE_TRADE_COMMAND_ADMIN_ADJUSTMENT:
		return "Admin Manual Balance Adjustment"
	case BALANCE_FREEZE_TRADE_COMMAND_MERCHANT_CALLBACK_FAILD:
		return "Callback Failed"
	default:
		return ""
	}
}

func BalanceSettlementTradeCommandMapping(status int64) string {
	switch status {
	case BALANCE_SETTLEMENT_TRADE_COMMAND_PAYIN:
		return "Unsettled"
	case BALANCE_SETTLEMENT_TRADE_COMMAND_TRANSACTION:
		return "Settled"
	default:
		return ""
	}
}

func IsReleaseMapping(isRelease int64) string {
	switch isRelease {
	case IS_RELEASE_NO:
		return "Unsettled"
	case IS_RELEASE_YES:
		return "Settled"
	default:
		return ""
	}
}

func OrderStatusMapping(orderStatus int64) string {
	switch orderStatus {
	case AUDIT_STATUS_WAIT_AUDIT:
		return "Under merchant review"
	case AUDIT_STATUS_PASS:
		return "Approved"
	case AUDIT_STATUS_REJECT:
		return "Merchant rejection"
	case AUDIT_STATUS_ADMIN_WAIT_AUDIT:
		return "Under platform review"
	case AUDIT_STATUS_ADMIN_REJECT:
		return "Platform rejection"
	case AUDIT_STATUS_PASS_NOT_PAYMENT:
		return "Awaiting payment"
	default:
		return ""
	}
}

func CommonStatusMapping(status int64) string {
	switch status {
	case STATUS_ENABLE:
		return "On"
	case STATUS_DISABLE:
		return "Off"
	default:
		return ""
	}
}

func DepositTypeMapping(depositType int64) string {
	switch depositType {
	case DEPOSIT_ORDER_TYPE_APP:
		return "APP"
	case DEPOSIT_ORDER_TYPE_BANK:
		return "Bank"
	case DEPOSIT_ORDER_TYPE_BORROW:
		return "Borrow First, Repay Later"
	default:
		return ""
	}
}

func UndertakeMapping(undertake int64) string {
	switch undertake {
	case UNDERTAKE_PLATFORM:
		return "Platform"
	case UNDERTAKE_MERCHANT:
		return "Merchant"
	case UNDERTAKE_BANK:
		return "Bank"
	default:
		return ""
	}
}

func RefundTypeMapping(refundType int64) string {
	switch refundType {
	case REFUND_TYPE_REFUND:
		return "Reversed"
	case REFUND_TYPE_COMPLAIN:
		return "Complaint"
	default:
		return ""
	}
}

func OrderTypeMapping(orderType int64) string {
	switch orderType {
	case ORDER_TYPE_PAY_IN:
		return "Pay in"
	case ORDER_TYPE_PAY_OUT:
		return "Pay out"
	default:
		return ""
	}
}

func AmountTypeCommandMapping(amountType int64) string {
	switch amountType {
	case BALANCE_AMOUNT_CHANGE:
		return "Transaction Amount"
	case BALANCE_FEE_CHANGE:
		return "Rate Fee"
	default:
		return ""
	}
}

func DepositOrderTransactionStatusMapping(transactionStatus int64) string {
	switch transactionStatus {
	case DEPOSIT_ORDER_TRANSFER_STATUS_PENDING:
		return "Pending"
	case DEPOSIT_ORDER_TRANSFER_STATUS_SUCCESS:
		return "Success"
	case DEPOSIT_ORDER_TRANSFER_STATUS_FAILED:
		return "Failure"
	case DEPOSIT_ORDER_TRANSFER_STATUS_REVERSAL:
		return "Refunded"
	default:
		return ""
	}
}
