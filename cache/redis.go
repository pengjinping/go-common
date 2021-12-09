package cache

import (
	"fmt"
	"git.kuainiujinke.com/oa/oa-common/config"
	timeHelper "git.kuainiujinke.com/oa/oa-common/utils/time"
	"github.com/garyburd/redigo/redis"
	"github.com/techoner/gophp/serialize"
	"time"
)

type RedisStore struct {
	pool              *redis.Pool
	defaultExpiration time.Duration
}

func NewRedisStore() *RedisStore {
	var conf config.RedisConfig
	if err := config.UnmarshalKey("Redis", &conf); err != nil {
		fmt.Printf("Redis config init failed: %v\n", err)
	}
	address := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	fmt.Printf("%s Redis连接信息: %s - ", timeHelper.FormatDateTime(time.Now()), address)

	pool := &redis.Pool{
		MaxActive:   512,
		MaxIdle:     10,
		Wait:        false,
		IdleTimeout: 3 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				panic(fmt.Sprintf("连接失败：%v", err))
			}

			if len(conf.Password) > 0 { // 有密码的情况
				if _, err := c.Do("AUTH", conf.Password); err != nil {
					fmt.Printf("密码错误：%v", err)
					err := c.Close()
					if err != nil {
						return nil, err
					}
					return nil, err
				}
			} else {
				if _, err := c.Do("ping"); err != nil {
					fmt.Printf("请求Ping错误：%v", err)
					err := c.Close()
					if err != nil {
						return nil, err
					}
					return nil, err
				}
			}

			//分库，默认是0，使用SELECT key，可进行切换，select 0切换到默认数据库
			if conf.DBName > 0 {
				if _, err := c.Do("SELECT", conf.DBName); err != nil {
					fmt.Printf("选择分库DB[%d]失败：%v", conf.DBName, err)
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

	//测试连接
	conn := pool.Get()
	if err := conn.Err(); err != nil {
		return nil
	}
	fmt.Println("success")

	return &RedisStore{pool: pool, defaultExpiration: conf.Expiration}
}

func (c *RedisStore) Set(key string, value interface{}, t int) {
	conn := c.pool.Get()
	defer conn.Close()

	out, _ := serialize.Marshal(value) //序列化操作，序列化可以保存对象
	if t > 0 {
		conn.Do("setex", key, t, string(out))
	} else {
		conn.Do("set", key, string(out))
	}
}

func (c *RedisStore) Forever(key string, value interface{}) {
	conn := c.pool.Get()
	defer conn.Close()

	out, _ := serialize.Marshal(value) //序列化操作，序列化可以保存对象
	conn.Do("set", key, string(out))
}

func (c *RedisStore) Get(key string) interface{} {
	conn := c.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("get", key)
	if err != nil || reply == nil {
		return nil
	}
	out, _ := redis.Bytes(reply, err)
	res, _ := serialize.UnMarshal(out)
	return res
}

func (c *RedisStore) Delete(key string) {
	conn := c.pool.Get()
	defer conn.Close()

	conn.Do("del", key)
}

func (c *RedisStore) IsExpire(key string) bool {
	conn := c.pool.Get()
	defer conn.Close()

	b, err := redis.Bool(conn.Do("exists", key))
	if b || err != nil {
		return false
	}

	return true
}

func (c *RedisStore) Has(key string) bool {
	return !c.IsExpire(key)
}
