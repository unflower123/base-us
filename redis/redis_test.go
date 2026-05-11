package redis

import (
	"context"
	"testing"
)

func Test_ReidsDeductBalance(t *testing.T) {
	conf := RedisConf{
		Name:       "redis",
		Addr:       "127.0.0.1:6379",
		Password:   "",
		DB:         0,
		UseCluster: false,
	}

	redis := NewRedis(&conf)

	key := "merchant:balance:1871425725003008"
	type balanceInfo struct {
		Id              int
		Balance         int
		WithdrawBalance int
		FreezeBalance   int
		UnSettleBalance int
	}
	info := balanceInfo{
		Id:              1871425725003008,
		Balance:         10000,
		WithdrawBalance: 1000,
		FreezeBalance:   6000,
		UnSettleBalance: 4000,
	}
	redis.HMSet(context.Background(), key, info)

}
