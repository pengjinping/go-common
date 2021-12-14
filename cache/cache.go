package cache

import (
	"context"
	"fmt"
	"log"
	"strings"

	"git.kuainiujinke.com/oa/oa-go-common/config"
)

var Stores = make(map[string]StoreInterface)

type Cache struct {
	driver string // 缓存驱动
	store  StoreInterface
}

type StoreInterface interface {
	GetTenant() string                           // 获取租户信息
	SetTenant(tenant string, tenantId int) bool  // 设置连接库DB
	Set(key string, value interface{}, time int) // 设置缓存带过期时间
	Forever(key string, value interface{})       // 设置永久缓存无过期时间
	Get(key string) interface{}                  // 获取缓存数据
	Delete(key string)                           // 删除key
	Has(key string) bool                         // 判断key是否存在
	IsExpire(key string) bool                    // 判断key是否过期
	GC()                                         // 随机删除已过期key
	Keys() interface{}                           // 获取所有key
}

func Init() {
	// 需要使用的时候直接获取，没有注入到全局变量中
	driver := getConfigCacheDriver()
	register(driver)
}

func register(driver string) {
	var store StoreInterface
	if driver == "redis" {
		store = NewRedisStore(config.SystemTenant)
	} else if driver == "memory" {
		store = NewMemoryStore(config.SystemTenant)
		go store.GC() // 开启一个协成 清理过期缓存数据
	} else {
		log.Printf("缓存驱动 \"%s\" 不存在. 可选择的缓存驱动: memory, redis\n", driver)
		return
	}

	if _, ok := Stores[driver]; ok {
		return
	}

	Stores[driver] = store
	log.Printf("Cache driver \"%s\" connected success.", driver)
}

func getConfigCacheDriver() string {
	var conf config.CacheConfig
	if err := config.UnmarshalKey("Cache", &conf); err != nil {
		fmt.Printf("Cache config init failed: %v\n", err)
	}
	driver := conf.Driver
	if len(driver) <= 0 {
		driver = "memory"
	}

	return driver
}

func GetDefault(ctx context.Context) *Cache {
	driver := getConfigCacheDriver()
	return GetByDriver(ctx, driver)
}

func GetByDriver(ctx context.Context, driver string) *Cache {
	driver = strings.ToLower(driver)

	var ca *Cache
	ca, tenant, siteID := cacheTenant(ctx)
	if ca == nil || ca.driver != driver {
		if _, ok := Stores[driver]; !ok {
			register(driver)
		}

		ca = &Cache{
			driver: driver,
			store:  Stores[driver],
		}
	}

	ca.SetTenant(tenant, siteID)
	return ca
}

func cacheTenant(ctx context.Context) (*Cache, string, int) {
	var ca *Cache
	if caName := ctx.Value("cache"); caName != nil {
		ca = caName.(*Cache)
	}

	tenant := config.SystemTenant
	if tenantName := ctx.Value("tenant"); tenantName != nil {
		tenant = tenantName.(string)
	}

	siteID := 0
	if tenantIdName := ctx.Value("siteID"); tenantIdName != nil {
		siteID = tenantIdName.(int)
	}

	return ca, tenant, siteID
}

/** 结构体方法 */

func (c *Cache) GetStoreName() string {
	return c.driver
}
func (c *Cache) GetTenant() string {
	return c.store.GetTenant()
}
func (c *Cache) SetTenant(tenant string, tenantId int) bool {
	return c.store.SetTenant(tenant, tenantId)
}

func (c *Cache) Set(key string, value interface{}, time int) {
	c.store.Set(key, value, time)
}
func (c *Cache) Forever(key string, value interface{}) {
	c.store.Forever(key, value)
}
func (c *Cache) Get(key string) interface{} {
	return c.store.Get(key)
}

func (c *Cache) Delete(key string) {
	c.store.Delete(key)
}
func (c *Cache) Has(key string) bool {
	return c.store.Has(key)
}
func (c *Cache) IsExpire(key string) bool {
	return c.store.IsExpire(key)
}
func (c *Cache) Keys() interface{} {
	return c.store.Keys()
}

func (c *Cache) Remember(key string, time int, f func(...interface{}) (interface{}, error), args ...interface{}) interface{} {
	isExpire := c.store.IsExpire(key)
	if !isExpire {
		return c.Get(key)
	}

	value, err := f(args...)
	if err != nil {
		return nil
	}

	c.Set(key, value, time)

	return value
}
