/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/4/18 10:47
 */
package metadata

import (
	"base/consts"
	"context"
	"google.golang.org/grpc/metadata"
	"strconv"
)

func AddMerchantId(ctx context.Context, value string) (context.Context, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	md.Set(consts.METADATA_KEY_MERCHANT_ID, value)
	return metadata.NewOutgoingContext(ctx, md), nil
}

func AddUserId(ctx context.Context, value string) (context.Context, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	md.Set(consts.METADATA_KEY_USER_ID, value)

	return metadata.NewOutgoingContext(ctx, md), nil
}

func AddUserName(ctx context.Context, value string) (context.Context, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	md.Set(consts.METADATA_KEY_USER_NAME, value)

	return metadata.NewOutgoingContext(ctx, md), nil
}

func AddRoleId(ctx context.Context, value string) (context.Context, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	md.Set(consts.METADATA_KEY_ROLE_ID, value)

	return metadata.NewOutgoingContext(ctx, md), nil
}

// SetMetadata creates a new outgoing context with the provided metadata from a map.
// It replaces any existing metadata in the context.
func SetMetadata(ctx context.Context, data map[string]string) (context.Context, error) {
	md := metadata.New(nil)
	for key, value := range data {
		md.Set(key, value)
	}
	return metadata.NewOutgoingContext(ctx, md), nil
}

// AppendMetadata appends multiple metadata key-value pairs from a map
// to the existing outgoing context.
func AppendMetadata(ctx context.Context, data map[string]string) (context.Context, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	for key, value := range data {
		md.Append(key, value) // Use Append to add to existing values
	}

	return metadata.NewOutgoingContext(ctx, md), nil
}

func SetMetadataMap(ctx context.Context) (context.Context, error) {
	mapData := map[string]string{
		consts.METADATA_KEY_MERCHANT_ID: "1870077057774848",
		consts.METADATA_KEY_APPID:       "11afbf373e0b4777f4d56522289fde57",
		consts.METADATA_KEY_USER_ID:     "112233",
		consts.METADATA_KEY_USER_NAME:   "lishi",
		consts.METADATA_KEY_ROLE_ID:     "1",
	}
	return SetMetadata(ctx, mapData)
}

func SetMetadataInfo(ctx context.Context, metadata Metadata) (context.Context, error) {
	mapData := map[string]string{
		consts.METADATA_KEY_MERCHANT_ID:   strconv.FormatInt(metadata.MerchantId, 10),
		consts.METADATA_KEY_APPID:         metadata.MerchantAppid,
		consts.METADATA_KEY_USER_ID:       strconv.FormatInt(metadata.UserId, 10),
		consts.METADATA_KEY_USER_NAME:     metadata.UserName,
		consts.METADATA_KEY_ROLE_ID:       strconv.FormatInt(metadata.RoleId, 10),
		consts.METADATA_KEY_CURRENCY:      metadata.Currency,
		consts.METADATA_KEY_MERCHANT_NAME: metadata.MerchantName,
	}
	return SetMetadata(ctx, mapData)
}
