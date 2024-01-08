package rds

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	DefaultExpiry = 8 * time.Second
	DefaultTries  = 1
	DefaultDelay  = 500 * time.Millisecond
	DefaultFactor = 0.01
)

var (
	ErrMaxTries = errors.New("reached max tries")
)

type RedisMutex struct {
	Key    string
	Expiry time.Duration

	Tries int
	Delay time.Duration

	Factor float64

	value string
	until time.Time

	redis redis.UniversalClient
}

type MutexOptFunc func(m *RedisMutex)

func WithExpiry(expiry time.Duration) MutexOptFunc {
	return func(m *RedisMutex) {
		m.Expiry = expiry
	}
}

func WithTries(tries int) MutexOptFunc {
	return func(m *RedisMutex) {
		m.Tries = tries
	}
}

func WithDelay(delay time.Duration) MutexOptFunc {
	return func(m *RedisMutex) {
		m.Delay = delay
	}
}

func WithFactor(factor float64) MutexOptFunc {
	return func(m *RedisMutex) {
		m.Factor = factor
	}
}

func NewRedisMutex(key string, redisClient redis.UniversalClient, opts ...MutexOptFunc) *RedisMutex {
	r := &RedisMutex{
		Key:    key,
		Expiry: DefaultExpiry,
		Tries:  DefaultTries,
		Delay:  DefaultDelay,
		Factor: DefaultFactor,
		redis:  redisClient,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (m *RedisMutex) LockKeeper() (unlocked <-chan struct{}, err error) {
	err = m.Lock()
	if err != nil {
		return
	}

	sigCh := make(chan struct{})
	go func() {
		for {
			time.Sleep(m.Expiry / 2)
			if m.Touch() {
				continue
			}
			close(sigCh)
			return
		}
	}()
	unlocked = sigCh
	return
}

func (m *RedisMutex) Lock() error {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	value := base64.StdEncoding.EncodeToString(b)

	expiry := m.Expiry
	if expiry == 0 {
		expiry = DefaultExpiry
	}

	retries := m.Tries
	if retries == 0 {
		retries = DefaultTries
	}

	delay := m.Delay
	if delay == 0 {
		delay = DefaultDelay
	}

	for i := 0; i < retries; i++ {
		ok, err := m.tryLock(m.Key, value, expiry)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}

		if i == retries-1 {
			break
		}

		time.Sleep(delay)
	}

	return ErrMaxTries
}

func (m *RedisMutex) tryLock(key, value string, expiry time.Duration) (bool, error) {
	start := time.Now()

	var ok bool
	boolCmd := m.redis.SetNX(context.Background(), key, value, expiry)
	if boolCmd.Err() != nil {
		return false, boolCmd.Err()
	}
	ok = boolCmd.Val()

	factor := m.Factor
	if factor == 0 {
		factor = DefaultFactor
	}

	until := time.Now().Add(expiry - time.Since(start) - time.Duration(int64(float64(expiry)*factor)) + 2*time.Millisecond)
	if ok && time.Now().Before(until) {
		m.value = value
		m.until = until
		return true, nil
	}

	cmd := redis.NewScript(delScript).Run(context.Background(), m.redis, []string{key}, m.value)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	return false, nil
}

func (m *RedisMutex) Touch() bool {
	value := m.value
	if value == "" {
		return false
	}

	expiry := m.Expiry
	if expiry == 0 {
		expiry = DefaultExpiry
	}
	cmd := redis.NewScript(touchScript).Run(context.Background(), m.redis, []string{m.Key}, m.value, expiry)
	return cmd.Err() == nil
}

func (m *RedisMutex) Unlock() bool {
	value := m.value
	if value == "" {
		// panic("redis mutex: unlock of unlocked mutex")
		return false
	}

	m.value = ""
	m.until = time.Unix(0, 0)

	cmd := redis.NewScript(delScript).Run(context.Background(), m.redis, []string{m.Key}, value)
	return cmd.Err() == nil
}

var delScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end`

var touchScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("set", KEYS[1], ARGV[1], "xx", "px", ARGV[2])
else
	return "ERR"
end`
