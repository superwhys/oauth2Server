package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/superwhys/superGo/superLog"
	"time"
)

type Store struct {
	pool    *redis.Pool
	timeout time.Duration
}

// DialRedisPoolBlocked 当连接池conn用完时，会阻塞等待获取conn
func DialRedisPoolBlocked(addr string, db int, maxIdle int, timeout time.Duration, password ...string) *Store {
	pool := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxIdle,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial:        redisDialFunc(addr, db, password...),
	}
	return &Store{
		pool:    pool,
		timeout: timeout,
	}
}

func (s *Store) Close() {
	s.pool.Close()
}

func redisDialFunc(addr string, db int, password ...string) func() (redis.Conn, error) {
	return func() (redis.Conn, error) {

		superLog.Debug("Dial redis", addr, db, addr)
		options := []redis.DialOption{
			redis.DialDatabase(db),
			redis.DialConnectTimeout(5 * time.Second),
		}
		if len(password) > 0 {
			options = append(options, redis.DialPassword(password[0]))
		}
		return redis.Dial("tcp", addr, options...)
	}
}

func (s *Store) Get(key string) (string, error) {
	conn := s.pool.Get()
	defer conn.Close()

	if resp, err := conn.Do("GET", key); err != nil {
		return "", err
	} else {
		var text string
		if resp == nil {
			text = ""
		} else {
			text = string(resp.([]byte))
		}
		return text, nil
	}
}

func (s *Store) Set(key, value string) error {
	conn := s.pool.Get()
	defer conn.Close()

	if _, err := conn.Do("SET", key, value, "EX", s.timeout.Seconds(), "NX"); err != nil {
		return err
	}
	return nil
}
