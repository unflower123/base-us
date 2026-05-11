/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/4/18 10:47
 */
package metadata

import (
	"context"
	"strconv"

	"base/consts"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

type Metadata struct {
	MerchantId    int64  `json:"merchant_id"`
	UserId        int64  `json:"user_id"`
	RoleId        int64  `json:"role_id"`
	MerchantAppid string `json:"merchant_appid"`
	UserName      string `json:"user_name"`
	Email         string `json:"email"`
	AuthSecret    string `json:"auth_secret"`
	Currency      string `json:"currency"`
	MerchantName  string `json:"merchant_name"`
}

// ExtractMerchantId extracts merchant_id from the context metadata
func ExtractMerchantId(ctx context.Context) (int64, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, ok := md[consts.METADATA_KEY_MERCHANT_ID]; ok && len(val) > 0 {
			merchantId, err := strconv.ParseInt(val[0], 10, 64)
			if err != nil {
				return 0, errors.Wrap(err, "failed to parse merchant id from metadata")
			}
			return merchantId, nil
		}
		return 0, errors.New("merchant id not found in metadata")
	}
	return 0, errors.New("metadata not found in context")
}

func ExtractMerchantAppid(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, ok := md[consts.METADATA_KEY_APPID]; ok && len(val) > 0 {
			return val[0], nil
		}
		return "", errors.New("merchant appid not found in metadata")
	}
	return "", errors.New("metadata not found in context")
}

func ExtractUserId(ctx context.Context) (int64, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, ok := md[consts.METADATA_KEY_USER_ID]; ok && len(val) > 0 {
			userId, err := strconv.ParseInt(val[0], 10, 64)
			if err != nil {
				return 0, errors.Wrap(err, "failed to parse user id from metadata")
			}
			return userId, nil
		}
		return 0, errors.New("user id not found in metadata")
	}
	return 0, errors.New("metadata not found in context")
}

func ExtractUserName(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, ok := md[consts.METADATA_KEY_USER_NAME]; ok && len(val) > 0 {
			if len(val[0]) <= 0 {
				return "", errors.New("failed user name from metadata null value")
			}
			return val[0], nil
		}
		return "", errors.New("user name not found in metadata")
	}
	return "", errors.New("metadata not found in context")
}

func ExtractRoleId(ctx context.Context) (int64, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, ok := md[consts.METADATA_KEY_ROLE_ID]; ok && len(val) > 0 {
			roleId, err := strconv.ParseInt(val[0], 10, 64)
			if err != nil {
				return 0, errors.Wrap(err, "failed to parse role id from metadata")
			}
			return roleId, nil
		}
		return 0, errors.New("role id not found in metadata")
	}
	return 0, errors.New("metadata not found in context")
}

func ExtractCurrency(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, ok := md[consts.METADATA_KEY_CURRENCY]; ok && len(val) > 0 {
			if len(val[0]) <= 0 {
				return "", errors.New("failed currency from metadata null value")
			}
			return val[0], nil
		}
		return "", errors.New("currency not found in metadata")
	}
	return "", errors.New("metadata not found in context")
}

func ExtractMerchantName(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, ok := md[consts.METADATA_KEY_MERCHANT_NAME]; ok && len(val) > 0 {
			if len(val[0]) <= 0 {
				return "", errors.New("failed merchant name from metadata null value")
			}
			return val[0], nil
		}
		return "", errors.New("merchant name not found in metadata")
	}
	return "", errors.New("metadata not found in context")
}

func ExtractMetadata(ctx context.Context) (Metadata, error) {
	var md Metadata
	var err error

	md.MerchantId, err = ExtractMerchantId(ctx)
	if err != nil {
		return md, err
	}

	md.MerchantAppid, err = ExtractMerchantAppid(ctx)
	if err != nil {
		return md, err
	}

	md.UserId, err = ExtractUserId(ctx)
	if err != nil {
		return md, err
	}

	md.UserName, err = ExtractUserName(ctx)
	if err != nil {
		return md, err
	}

	md.RoleId, err = ExtractRoleId(ctx)
	if err != nil {
		return md, err
	}

	md.Currency, err = ExtractCurrency(ctx)
	if err != nil {
		return md, err
	}

	md.MerchantName, err = ExtractMerchantName(ctx)
	if err != nil {
		return md, err
	}

	return md, nil
}
