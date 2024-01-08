package rds

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"

	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
)

const (
	DefaultTimeout = time.Minute
	ZeroTimeOut    = 0
)

var (
	ErrMissedKey = errors.New("err_missed_key")
)

func getTimeoutDur(params ...int64) time.Duration {
	var timeout time.Duration
	if len(params) > 0 && params[0] > 0 {
		timeout = time.Duration(params[0]) * time.Second
		return timeout
	}
	return ZeroTimeOut
}

type RedisConfig struct {
	Addrs    []string
	Password string
	PoolSize int
	DB       int
	// The sentinel master name.
	// Only failover clients.
	MasterName       string
	SentinelPassword string
	Serializer       ISerializer
	KeyPrefix        string
}

type RedisClient struct {
	redis.UniversalClient
	*RedisConfig
}

func NewRedisProvider(config *RedisConfig) (client *RedisClient, err error) {
	if config.Serializer == nil {
		config.Serializer = &gobSerializer{}
	}
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	opt := redis.UniversalOptions{}
	if err = json.Unmarshal(jsonBytes, &opt); err != nil {
		return nil, errorx.WithStack(err)
	}
	c := &RedisClient{
		redis.NewUniversalClient(&opt),
		config,
	}

	return c, nil
}

func (c *RedisClient) SetSerializer(s ISerializer) {
	c.Serializer = s
}

func (p *RedisClient) Get(key string, value interface{}) error {
	key = p.KeyPrefix + key
	stringCmd := p.UniversalClient.Get(context.Background(), key)
	if stringCmd.Err() != nil {
		return ErrMissedKey
	}
	err := p.RedisConfig.Serializer.Unmarshal([]byte(stringCmd.Val()), value)
	if err != nil {
		return err
	}
	return nil
}

func (p *RedisClient) Set(key string, val interface{}, params ...int64) error {
	byt, err := p.RedisConfig.Serializer.Marshal(val)
	if err != nil {
		return err
	}
	key = p.KeyPrefix + key
	timeout := getTimeoutDur(params...)
	cmd := p.UniversalClient.Set(context.Background(), key, string(byt), timeout)
	return cmd.Err()
}

func (p *RedisClient) SetNx(key string, val interface{}, params ...int64) error {
	byt, err := p.RedisConfig.Serializer.Marshal(val)
	if err != nil {
		return err
	}
	key = p.KeyPrefix + key
	timeout := getTimeoutDur(params...)

	// args := []interface{}{"set", key, string(byt), "NX"}
	// if timeout > 0 {
	//	args = append(args, "EX", strconv.FormatFloat(timeout.Seconds(), 'f', -1, 64))
	// }
	cmd := p.UniversalClient.SetNX(context.Background(), key, string(byt), timeout)
	return cmd.Err()
}

func (p *RedisClient) Touch(key string, params ...int64) (err error) {
	key = p.KeyPrefix + key
	timeout := getTimeoutDur(params...)

	bCmd := p.UniversalClient.Expire(context.Background(), key, timeout)
	if bCmd.Err() != nil {
		return bCmd.Err()
	}
	if !bCmd.Val() {
		return ErrMissedKey
	}
	return
}

func (p *RedisClient) Delete(key string) (err error) {
	key = p.KeyPrefix + key

	delCmd := p.UniversalClient.Del(context.Background(), key)
	if delCmd.Err() != nil {
		return delCmd.Err()
	}
	if delCmd.Val() == 0 {
		err = ErrMissedKey
	}
	return
}

func (p *RedisClient) Incr(key string, params ...int64) (res int64, err error) {
	key = p.KeyPrefix + key
	var cnt int64 = 1
	if len(params) > 0 {
		cnt = params[0]
	}
	intCmd := p.UniversalClient.IncrBy(context.Background(), key, cnt)
	if intCmd.Err() != nil {
		err = intCmd.Err()
		return
	}
	return intCmd.Val(), nil
}

func (p *RedisClient) Decr(key string, params ...int64) (res int64, err error) {
	key = p.KeyPrefix + key
	var cnt int64 = 1
	if len(params) > 0 {
		cnt = params[0]
	}
	intCmd := p.UniversalClient.DecrBy(context.Background(), key, cnt)
	if intCmd.Err() != nil {
		err = intCmd.Err()
		return
	}
	return intCmd.Val(), nil
}

func (p *RedisClient) Exists(key string) (exists bool, err error) {
	key = p.KeyPrefix + key

	intCmd := p.UniversalClient.Exists(context.Background(), key)
	if intCmd.Err() != nil {
		err = intCmd.Err()
		return
	}
	exists = intCmd.Val() == 1
	return
}

func (p *RedisClient) GC() error {
	return nil
}

func (p *RedisClient) TTL(key string) (time.Duration, error) {
	durationCmd := p.UniversalClient.TTL(context.Background(), key)
	if durationCmd.Err() != nil {
		return 0, durationCmd.Err()
	}
	return durationCmd.Val(), nil
}
