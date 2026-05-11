package metadata

import (
	"base/consts"
	"context"
	"errors"
	"strconv"
	"strings"
)

type AgentMetadata struct {
	UserID        uint64   `json:"user_id"`
	UserName      string   `json:"user_name"`
	UserAccount   string   `json:"user_account"`
	MerchantIDStr string   `json:"merchant_id_str"`
	BankIDStr     string   `json:"bank_id_str"`
	MerchantIDs   []uint64 `json:"merchant_ids"`
	BankIDs       []uint64 `json:"bank_ids"`
}

func SetAgentUserInfo(ctx context.Context, metadata *AgentMetadata) context.Context {
	bankIds := make([]uint64, 0)
	merchantIds := make([]uint64, 0)
	seenBankIds := make(map[uint64]struct{})
	seenMerchantIds := make(map[uint64]struct{})
	if metadata.MerchantIDStr != "" && metadata.MerchantIDStr != "-1" {
		strIds := strings.Split(metadata.MerchantIDStr, ",")
		for _, strId := range strIds {
			trimmedStr := strings.TrimSpace(strId)
			if trimmedStr != "" {
				if id, err := strconv.ParseUint(trimmedStr, 10, 64); err == nil {
					if _, exists := seenMerchantIds[id]; !exists {
						seenMerchantIds[id] = struct{}{}
						merchantIds = append(merchantIds, id)
					}
				}
			}
		}
	}
	if metadata.BankIDStr != "-1" {
		strIds := strings.Split(metadata.BankIDStr, ",")
		for _, strId := range strIds {
			trimmedStr := strings.TrimSpace(strId)
			if trimmedStr != "" {
				if id, err := strconv.ParseUint(trimmedStr, 10, 64); err == nil {
					if _, exists := seenBankIds[id]; !exists {
						seenBankIds[id] = struct{}{}
						bankIds = append(bankIds, id)
					}
				}
			}
		}
	}
	metadata.MerchantIDs = merchantIds
	metadata.BankIDs = bankIds
	return context.WithValue(ctx, consts.AGENT_USER_INFO, metadata)
}

func GetAgentUserInfo(ctx context.Context) (metadata *AgentMetadata, err error) {
	val, ok := ctx.Value(consts.AGENT_USER_INFO).(*AgentMetadata)
	if ok {
		return val, nil
	}
	return nil, errors.New("ctx not exist user info")
}
