/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/4/22 15:05
 */
package metadata

import (
	"base/consts"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestExtractMerchantId(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		expectedID  int64
		expectedErr error
	}{
		{
			name: "Valid Merchant ID",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				consts.METADATA_KEY_MERCHANT_ID, "12345",
			)),
			expectedID:  12345,
			expectedErr: nil,
		},
		{
			name: "Invalid Merchant ID",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				consts.METADATA_KEY_MERCHANT_ID, "abc",
			)),
			expectedID:  0,
			expectedErr: errors.New("failed to parse merchant id from metadata"),
		},
		{
			name:        "Merchant ID Not Found",
			ctx:         context.Background(),
			expectedID:  0,
			expectedErr: errors.New("merchant id not found in metadata"),
		},
		{
			name: "Empty Merchant ID",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				consts.METADATA_KEY_MERCHANT_ID, "",
			)),
			expectedID:  0,
			expectedErr: errors.New("metadata not found in context"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ExtractMerchantId(tt.ctx)
			if tt.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}

func TestExtractUserId(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		expectedID  int64
		expectedErr error
	}{
		{
			name: "Valid User ID",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				consts.METADATA_KEY_USER_ID, "67890",
			)),
			expectedID:  67890,
			expectedErr: nil,
		},
		{
			name: "Invalid User ID",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				consts.METADATA_KEY_USER_ID, "xyz",
			)),
			expectedID:  0,
			expectedErr: errors.New("failed to parse user id from metadata"),
		},
		{
			name:        "User ID Not Found",
			ctx:         context.Background(),
			expectedID:  0,
			expectedErr: errors.New("user id not found in metadata"),
		},
		{
			name: "Empty User ID",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				consts.METADATA_KEY_USER_ID, "",
			)),
			expectedID:  0,
			expectedErr: errors.New("user id not found in metadata"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ExtractUserId(tt.ctx)
			if tt.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}

func TestExtractUserName(t *testing.T) {
	tests := []struct {
		name         string
		ctx          context.Context
		expectedName string
		expectedErr  error
	}{
		{
			name: "Valid User Name",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				consts.METADATA_KEY_USER_NAME, "john_doe",
			)),
			expectedName: "john_doe",
			expectedErr:  nil,
		},
		{
			name:         "User Name Not Found",
			ctx:          context.Background(),
			expectedName: "",
			expectedErr:  errors.New("user name not found in metadata"),
		},
		{
			name: "Empty User Name",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				consts.METADATA_KEY_USER_NAME, "",
			)),
			expectedName: "",
			expectedErr:  errors.New("user name not found in metadata"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, err := ExtractUserName(tt.ctx)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				//assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedName, name)
			}
		})
	}
}

func TestExtractRoleId(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		expectedID  int64
		expectedErr error
	}{
		{
			name: "Valid Role ID",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				consts.METADATA_KEY_ROLE_ID, "101",
			)),
			expectedID:  101,
			expectedErr: nil,
		},
		{
			name: "Invalid Role ID",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				consts.METADATA_KEY_ROLE_ID, "abc",
			)),
			expectedID:  0,
			expectedErr: errors.New("failed to parse role id from metadata"),
		},
		{
			name:        "Role ID Not Found",
			ctx:         context.Background(),
			expectedID:  0,
			expectedErr: errors.New("role id not found in metadata"),
		},
		{
			name: "Empty Role ID",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(
				consts.METADATA_KEY_ROLE_ID, "",
			)),
			expectedID:  0,
			expectedErr: errors.New("role id not found in metadata"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ExtractRoleId(tt.ctx)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				//assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}
