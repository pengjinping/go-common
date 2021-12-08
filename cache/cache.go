package cache

import (
	"fmt"
	"log"
	"oa-common/config"
	"strings"
)

var Stores = make(map[string]StoreInterface)
var tenantPrefix string

type Cache struct {
	name  string
	store StoreInterface
}

type StoreInterface interface {
	GetStoreName() string                              // 获取缓存驱动
	Get(key string) (interface{}, error)               // 获取缓存数据
	Set(key string, value interface{}, time int) error // 设置缓存带过期时间
	Forever(key string, value interface{}) error       // 设置永久缓存无过期时间
	Delete(key string) error                           // 删除key
	Has(key string) (bool, error)                      // 判断key是否存在
}

func Init() {
	// 需要使用的时候直接获取，没有注入到全局变量中
	driver := getConfigCacheDriver()
	if driver == "redis" {
		register(RedisName, NewRedisStore())
	} else {
		register(MemoryName, NewMemoryStore())
	}
}
func register(name string, store StoreInterface) {
	if store == nil {
		log.Panic("Cache: Register store is nil")
	}

	name = strings.ToLower(name)
	if _, ok := Stores[name]; ok {
		log.Panic("Cache: Register store is exist")
	}

	Stores[name] = store
	log.Printf("cache \"%s\" connected success. 可选择的缓存驱动:memory, redis", name)
}

func getConfigCacheDriver() string {
	conf := config.Get("Cache").(config.CacheConfig)
	return strings.ToLower(conf.Driver)
}

func SetPrefix(prefix string) {
	tenantPrefix = prefix + ":"
}

func Get() (*Cache, error) {
	driver := getConfigCacheDriver()
	if len(driver) <= 0 {
		driver = "memory"
	}

	if store, ok := Stores[driver]; ok {
		return &Cache{
			name:  driver, //有点多余
			store: store,
		}, nil
	}

	return nil, fmt.Errorf("缓存驱动未设置")
}

/** 结构体方法 */

func (c *Cache) GetStoreName() string {
	return c.store.GetStoreName()
}

func (c *Cache) Set(key string, value interface{}, time int) error {
	return c.store.Set(tenantPrefix+key, value, time)
}
func (c *Cache) Forever(key string, value interface{}) error {
	return c.store.Forever(tenantPrefix+key, value)
}
func (c *Cache) Get(key string) (interface{}, error) {
	return c.store.Get(tenantPrefix + key)
}
func (c *Cache) Delete(key string) error {
	return c.store.Delete(tenantPrefix + key)
}
func (c *Cache) Has(key string) (bool, error) {
	return c.store.Has(tenantPrefix + key)
}
