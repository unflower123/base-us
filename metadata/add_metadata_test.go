/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/4/22 15:24
 */
package metadata

import (
	"base/consts"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestAddMerchantId(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		value       string
		expectedMD  metadata.MD
		expectedErr error
	}{
		{
			name:  "Add Merchant ID to Existing Metadata",
			ctx:   metadata.NewOutgoingContext(context.Background(), metadata.Pairs("key", "value")),
			value: "12345",
			expectedMD: metadata.MD{
				"key":                           []string{"value"},
				consts.METADATA_KEY_MERCHANT_ID: []string{"12345"},
			},
			expectedErr: nil,
		},
		{
			name:  "Add Merchant ID to New Metadata",
			ctx:   context.Background(),
			value: "12345",
			expectedMD: metadata.MD{
				consts.METADATA_KEY_MERCHANT_ID: []string{"12345"},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newCtx, err := AddMerchantId(tt.ctx, tt.value)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				//assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				md, ok := metadata.FromOutgoingContext(newCtx)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedMD, md)
			}
		})
	}
}

func TestAddUserId(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		value       string
		expectedMD  metadata.MD
		expectedErr error
	}{
		{
			name:  "Add User ID to Existing Metadata",
			ctx:   metadata.NewOutgoingContext(context.Background(), metadata.Pairs("key", "value")),
			value: "67890",
			expectedMD: metadata.MD{
				"key":                       []string{"value"},
				consts.METADATA_KEY_USER_ID: []string{"67890"},
			},
			expectedErr: nil,
		},
		{
			name:  "Add User ID to New Metadata",
			ctx:   context.Background(),
			value: "67890",
			expectedMD: metadata.MD{
				consts.METADATA_KEY_USER_ID: []string{"67890"},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newCtx, err := AddUserId(tt.ctx, tt.value)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				//assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				md, ok := metadata.FromOutgoingContext(newCtx)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedMD, md)
			}
		})
	}
}

func TestAddUserName(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		value       string
		expectedMD  metadata.MD
		expectedErr error
	}{
		{
			name:  "Add User Name to Existing Metadata",
			ctx:   metadata.NewOutgoingContext(context.Background(), metadata.Pairs("key", "value")),
			value: "john_doe",
			expectedMD: metadata.MD{
				"key":                         []string{"value"},
				consts.METADATA_KEY_USER_NAME: []string{"john_doe"},
			},
			expectedErr: nil,
		},
		{
			name:  "Add User Name to New Metadata",
			ctx:   context.Background(),
			value: "john_doe",
			expectedMD: metadata.MD{
				consts.METADATA_KEY_USER_NAME: []string{"john_doe"},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newCtx, err := AddUserName(tt.ctx, tt.value)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				//assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				md, ok := metadata.FromOutgoingContext(newCtx)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedMD, md)
			}
		})
	}
}

func TestAddRoleId(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		value       string
		expectedMD  metadata.MD
		expectedErr error
	}{
		{
			name:  "Add Role ID to Existing Metadata",
			ctx:   metadata.NewOutgoingContext(context.Background(), metadata.Pairs("key", "value")),
			value: "101",
			expectedMD: metadata.MD{
				"key":                       []string{"value"},
				consts.METADATA_KEY_ROLE_ID: []string{"101"},
			},
			expectedErr: nil,
		},
		{
			name:  "Add Role ID to New Metadata",
			ctx:   context.Background(),
			value: "101",
			expectedMD: metadata.MD{
				consts.METADATA_KEY_ROLE_ID: []string{"101"},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newCtx, err := AddRoleId(tt.ctx, tt.value)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				//assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				md, ok := metadata.FromOutgoingContext(newCtx)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedMD, md)
			}
		})
	}
}
