package metadata

import (
	"base/consts"
	"context"
	"errors"
)

func SetMerchantUserInfo(ctx context.Context, metadata *Metadata) context.Context {
	return context.WithValue(ctx, consts.MERCHANT_USER_INFO, metadata)
}

func GetMerchantUserInfo(ctx context.Context) (metadata *Metadata, err error) {
	val, ok := ctx.Value(consts.METADATA_KEY_MERCHANT_USER_INFO).(*Metadata)
	if ok {
		return val, nil
	}
	return nil, errors.New("ctx not exist merchant user info")
}

func GetMerchantID(ctx context.Context) (merchantID int64, err error) {
	val, ok := ctx.Value(consts.METADATA_KEY_MERCHANT_USER_INFO).(*Metadata)
	if ok {
		if val.MerchantId == 0 {
			return 0, errors.New("ctx data exception")
		}
		return val.MerchantId, nil
	}
	return 0, errors.New("ctx not exist merchant user info")
}

func GetMerchantAppID(ctx context.Context) (merchantAppID string, err error) {
	val, ok := ctx.Value(consts.METADATA_KEY_MERCHANT_USER_INFO).(*Metadata)
	if ok {
		if val.MerchantAppid == "" {
			return "", errors.New("ctx data exception")
		}
		return val.MerchantAppid, nil
	}
	return "", errors.New("ctx not exist merchant user info")
}
