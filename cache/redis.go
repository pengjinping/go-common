package cache

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/techoner/gophp/serialize"
	"oa-common/config"
	"time"
)

type RedisStore struct {
	pool              *redis.Pool
	defaultExpiration time.Duration
}

var RedisName = "Redis"

func NewRedisStore() *RedisStore {
	conf := config.Get("Redis").(config.RedisConfig)
	address := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

	pool := &redis.Pool{
		MaxActive:   512,
		MaxIdle:     10,
		Wait:        false,
		IdleTimeout: 3 * time.Second,
		Dial: func() (redis.Conn, error) {
			selectDb := redis.DialDatabase(conf.DBName)
			c, err := redis.Dial("tcp", address, selectDb)
			if err != nil {
				return nil, err
			}
			if len(conf.Password) > 0 { // 有密码的情况
				if _, err := c.Do("AUTH", conf.Password); err != nil {
					err := c.Close()
					if err != nil {
						return nil, err
					}
					return nil, err
				}
			} else { // 没有密码的时候 ping 连接
				if _, err := c.Do("ping"); err != nil {
					err := c.Close()
					if err != nil {
						return nil, err
					}
					return nil, err
				}
			}
			return c, err
		},
	}
	return &RedisStore{pool: pool, defaultExpiration: conf.Expiration}
}

func (c *RedisStore) GetStoreName() string {
	return RedisName
}

func (c *RedisStore) Set(key string, value interface{}, t int) error {
	conn := c.pool.Get()
	defer conn.Close()

	out, _ := serialize.Marshal(value) //序列化操作，序列化可以保存对象
	if t > 0 {
		_, err := conn.Do("setex", key, t, out)
		return err
	} else {
		 _, err := conn.Do("set", key, out)
		return err
	}
}

func (c *RedisStore) Forever(key string, value interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()

	out, _ := serialize.Marshal(value) //序列化操作，序列化可以保存对象
	_, err := conn.Do("set", key, out)
	return err
}

func (c *RedisStore) Get(key string) (interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("get", key)
	if err != nil {
		return nil, err
	}
	if reply == nil {
		return nil, fmt.Errorf("缓存不存在")
	}

	out, err := redis.Bytes(reply, err)
	return serialize.UnMarshal(out) //序列化操作，序列化可以保存对象
}

func (c *RedisStore) Delete(key string) error {
	conn := c.pool.Get()
	defer conn.Close()

	b, err := redis.Bool(conn.Do("exists", key))
	if err != nil {
		return err
	}
	if !b {
		return fmt.Errorf("缓存不存在")
	}
	_, errDel := conn.Do("del", key)
	if errDel != nil {
		return errDel
	}

	return nil
}

func (c *RedisStore) IsExpire(key string) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()

	b, err := redis.Bool(conn.Do("exists", key))
	if err != nil {
		return false, err
	}
	if b {
		return false, nil
	}

	return true, nil
}

func (c *RedisStore) Has(key string) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()

	b, err := redis.Bool(conn.Do("exists", key))
	if err != nil {
		return false, err
	}
	if !b {
		return false, nil
	}

	return true, nil
}

func (c *RedisStore) Clear() error {
	return nil
}

func (c *RedisStore) GetTTl(key string) (time.Duration, error) {
	return time.Duration(0), nil
}
