package redis

/**
 * @Author: guyu
 * @Desc:
 * @Date: 2025/4/16 14:55
 */
import (
	"context"
	"fmt"
	"github.com/fatih/structs"
	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

type RedisConf struct {
	Name       string
	Addr       string
	Password   string
	DB         int
	UseCluster bool
}

type Redis struct {
	Client *redis.Client
}

func NewRedis(conf *RedisConf) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Password,
		DB:           conf.DB,
		PoolSize:     200,
		MinIdleConns: 40,
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("redis connect ping error, err:%s", err.Error())
		return nil
	}

	log.Println("redis connect ping response: ", pong)
	return &Redis{Client: client}
}

func (r *Redis) GetRedisClient() *redis.Client {
	return r.Client
}

func (r *Redis) DeductBalance(ctx context.Context, key string, balancekey string, freezekey string, num int64) (incr int64, decr int64, err error) {
	balanceScript := redis.NewScript(`
	local key = KEYS[1]  
	local firstKey = KEYS[2]
	local secondKey = KEYS[3]	
	local change = tonumber(ARGV[1])
	
	local firstBalance = tonumber(redis.call('HGET', key, firstKey) or "0")
	local secondBalance = tonumber(redis.call('HGET', key,secondKey) or "0")
	if  change  > firstBalance then
	return {0}
	end
	
	local newFirstBalance = firstBalance - change
	local newSecondBalance = secondBalance + change
	redis.call('HSET', key,firstKey, newFirstBalance)
	redis.call('HSET', key,secondKey, newSecondBalance)
	return {1, newFirstBalance,newSecondBalance}
  `)

	result, err := balanceScript.Run(ctx, r.Client, []string{key, balancekey, freezekey}, num).Result()
	if err != nil {
		return
	}

	vals := result.([]interface{})
	if vals[0] == int64(0) {
		err = fmt.Errorf("balance is not enough")
		return
	}
	incr = vals[1].(int64)
	decr = vals[2].(int64)
	return
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	data, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			return "", err
		}
	}

	return data, nil
}

func (r *Redis) HMSet(ctx context.Context, key string, value interface{}) (bool, error) {
	m, ok := value.(map[string]interface{})
	if ok == false {
		m = structs.Map(value)
	}
	return r.Client.HMSet(ctx, key, m).Result()
}

// key prefix with "pb:"
func (r *Redis) SetMessage(ctx context.Context, key string, message proto.Message, expiration time.Duration) (err error) {

	data, err := proto.Marshal(message)
	if err != nil {
		//fmt.Printf("redis set proto marshal failed, err: %s", err.Error())
		return
	}

	err = r.Client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		//fmt.Printf("redisClient set error : %s", err)
		return
	}

	return
}

// key prefix with "pb:"
func (r *Redis) GetMessage(ctx context.Context, key string, message proto.Message) (err error) {

	bytes, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		//fmt.Printf("redisClient get failed,  error : %s", err)
		return
	}

	err = proto.Unmarshal([]byte(bytes), message)
	if err != nil {
		//fmt.Printf("redisClient get message, proto unmarshal failed,  error : %s", err)
		return
	}

	return
}

// key prefix with "pb:"
func (r *Redis) SetStruct(ctx context.Context, key string, value interface{}, message proto.Message, expiration time.Duration) (err error) {
	err = copier.Copy(message, value)
	if err != nil {
		//fmt.Printf("redis set  copy value failed, err:%s", err.Error())
		return
	}

	return r.SetMessage(ctx, key, message, expiration)
}

// key prefix with "pb:"
func (r *Redis) GetStruct(ctx context.Context, key string, value interface{}, message proto.Message) (err error) {
	err = r.GetMessage(ctx, key, message)
	if err != nil {
		return
	}
	err = copier.Copy(value, message)
	if err != nil {
		//fmt.Printf("redisClient get value copy failed,  error : %s", err)
		return
	}

	return
}

func (r *Redis) IncrBalance(ctx context.Context, key string, balancekey string, freezekey string, num int64) (incr int64, decr int64, err error) {
	balanceScript := redis.NewScript(`
	local key = KEYS[1]  
	local firstKey = KEYS[2]
	local secondKey = KEYS[3]	
	local change = tonumber(ARGV[1])
	
	local firstBalance = tonumber(redis.call('HGET', key, firstKey) or "0")
	local secondBalance = tonumber(redis.call('HGET', key,secondKey) or "0")
	local newFirstBalance = firstBalance + change
	local newSecondBalance = secondBalance + change
	redis.call('HSET', key,firstKey, newFirstBalance)
	redis.call('HSET', key,secondKey, newSecondBalance)
	return {1, newFirstBalance,newSecondBalance}
  `)

	result, err := balanceScript.Run(ctx, r.Client, []string{key, balancekey, freezekey}, num).Result()
	if err != nil {
		return
	}

	vals := result.([]interface{})

	incr = vals[1].(int64)
	decr = vals[2].(int64)
	return
}

func (r *Redis) DecrBalance(ctx context.Context, key string, balancekey string, freezekey string, num int64) (incr int64, decr int64, err error) {
	balanceScript := redis.NewScript(`
	local key = KEYS[1]  
	local firstKey = KEYS[2]
	local secondKey = KEYS[3]	
	local change = tonumber(ARGV[1])
	
	local firstBalance = tonumber(redis.call('HGET', key, firstKey) or "0")
	local secondBalance = tonumber(redis.call('HGET', key,secondKey) or "0")
	local newFirstBalance = firstBalance - change
	local newSecondBalance = secondBalance - change
	redis.call('HSET', key,firstKey, newFirstBalance)
	redis.call('HSET', key,secondKey, newSecondBalance)
	return {1, newFirstBalance,newSecondBalance}
  `)

	result, err := balanceScript.Run(ctx, r.Client, []string{key, balancekey, freezekey}, num).Result()
	if err != nil {
		return
	}

	vals := result.([]interface{})

	incr = vals[1].(int64)
	decr = vals[2].(int64)
	return
}

func (r *Redis) IncrWithdrawBalance(ctx context.Context, key string, balancekey string, freezekey string, num int64) (incr int64, decr int64, err error) {
	balanceScript := redis.NewScript(`
	local key = KEYS[1]  
	local firstKey = KEYS[2]
	local secondKey = KEYS[3]	
	local change = tonumber(ARGV[1])
	
	local firstBalance = tonumber(redis.call('HGET', key, firstKey) or "0")
	local secondBalance = tonumber(redis.call('HGET', key,secondKey) or "0")

	if  change  > secondBalance then
	return {0}
	end

	local newFirstBalance = firstBalance + change
	local newSecondBalance = secondBalance - change
	redis.call('HSET', key,firstKey, newFirstBalance)
	redis.call('HSET', key,secondKey, newSecondBalance)
	return {1, newFirstBalance,newSecondBalance}
  `)

	result, err := balanceScript.Run(ctx, r.Client, []string{key, balancekey, freezekey}, num).Result()
	if err != nil {
		return
	}

	vals := result.([]interface{})
	if vals[0] == int64(0) {
		err = fmt.Errorf("balance is not enough")
		return
	}
	incr = vals[1].(int64)
	decr = vals[2].(int64)
	return
}

// 增加1减少2 减少3减少4 (可用余额增加，冻结金额减小 &&  总金额减小 ，冻结金额减小)
func (r *Redis) IncrWithdrawDecrBalance(ctx context.Context, key string, balancekey string, freezekey string, num int64, dbalancekey string, dfreezekey string, dnum int64) (incr int64, decr int64, dincr int64, ddecr int64, err error) {
	balanceScript := redis.NewScript(`
	local key = KEYS[1]  
	local firstKey = KEYS[2]
	local secondKey = KEYS[3]	

	local dfirstKey = KEYS[4]
	local dsecondKey = KEYS[5]	

	local change = tonumber(ARGV[1])
    local dchange = tonumber(ARGV[2])

	local firstBalance = tonumber(redis.call('HGET', key, firstKey) or "0")	
    local dfirstBalance = tonumber(redis.call('HGET', key, dfirstKey) or "0")

	if  dchange  > dfirstBalance then
	return {0}
	end

	local secondBalance = tonumber(redis.call('HGET', key,secondKey) or "0")
	local newFirstBalance = firstBalance + change
	local newSecondBalance = secondBalance - change
	redis.call('HSET', key,firstKey, newFirstBalance)
	redis.call('HSET', key,secondKey, newSecondBalance)

	local dsecondBalance = tonumber(redis.call('HGET', key,dsecondKey) or "0")
	local dnewFirstBalance = dfirstBalance - dchange
	local dnewSecondBalance = dsecondBalance - dchange
	redis.call('HSET', key,dfirstKey, dnewFirstBalance)
	redis.call('HSET', key,dsecondKey, dnewSecondBalance)

	return {1, newFirstBalance,newSecondBalance,dnewFirstBalance,dnewSecondBalance}
  `)

	result, err := balanceScript.Run(ctx, r.Client, []string{key, balancekey, freezekey, dbalancekey, dfreezekey}, num, dnum).Result()
	if err != nil {
		return
	}

	vals := result.([]interface{})
	if vals[0] == int64(0) {
		err = fmt.Errorf("balance is not enough")
		return
	}
	incr = vals[1].(int64)
	decr = vals[2].(int64)
	dincr = vals[3].(int64)
	ddecr = vals[4].(int64)
	return
}

// 减少1增加2 减少3减少4 (可用余额减少，冻结金额增加 &&  总金额减小 ，冻结金额减小)
func (r *Redis) DecrWithdrawDecrBalance(ctx context.Context, key string, balancekey string, freezekey string, num int64, dbalancekey string, dfreezekey string, dnum int64) (incr int64, decr int64, dincr int64, ddecr int64, err error) {
	balanceScript := redis.NewScript(`
	local key = KEYS[1]  
	local firstKey = KEYS[2]
	local secondKey = KEYS[3]	

	local dfirstKey = KEYS[4]
	local dsecondKey = KEYS[5]	

	local change = tonumber(ARGV[1])
    local dchange = tonumber(ARGV[2])

	local firstBalance = tonumber(redis.call('HGET', key, firstKey) or "0")
    local dfirstBalance = tonumber(redis.call('HGET', key, dfirstKey) or "0")

	if  change  > firstBalance then
	return {0}
	end

	if  dchange  > dfirstBalance then
	return {0}
	end

	local secondBalance = tonumber(redis.call('HGET', key,secondKey) or "0")
	local newFirstBalance = firstBalance - change
	local newSecondBalance = secondBalance + change
	redis.call('HSET', key,firstKey, newFirstBalance)
	redis.call('HSET', key,secondKey, newSecondBalance)

	local dsecondBalance = tonumber(redis.call('HGET', key,dsecondKey) or "0")
	local dnewFirstBalance = dfirstBalance - dchange
	local dnewSecondBalance = dsecondBalance - dchange
	redis.call('HSET', key,dfirstKey, dnewFirstBalance)
	redis.call('HSET', key,dsecondKey, dnewSecondBalance)

	return {1, newFirstBalance,newSecondBalance,dnewFirstBalance,dnewSecondBalance}
  `)

	result, err := balanceScript.Run(ctx, r.Client, []string{key, balancekey, freezekey, dbalancekey, dfreezekey}, num, dnum).Result()
	if err != nil {
		return
	}

	vals := result.([]interface{})
	if vals[0] == int64(0) {
		err = fmt.Errorf("balance is not enough")
		return
	}
	incr = vals[1].(int64)
	decr = vals[2].(int64)
	dincr = vals[3].(int64)
	ddecr = vals[4].(int64)
	return
}
