package cache

import (
	"context"
	"fmt"
	"git.kuainiujinke.com/oa/oa-common-golang/config"
	"log"
	"strings"
)

var Stores = make(map[string]StoreInterface)

type Cache struct {
	driver string // 缓存驱动
	store  StoreInterface
}

type StoreInterface interface {
	Tenant() string                              // 获取租户UUID
	SetTenant(tenant string, tenantId uint) bool  // 设置连接库DB
	Set(key string, value interface{}, time int) // 设置缓存带过期时间
	Forever(key string, value interface{})       // 设置永久缓存无过期时间
	Get(key string) interface{}                  // 获取缓存数据
	Delete(key string)                           // 删除key
	Has(key string) bool                         // 判断key是否存在
	IsExpire(key string) bool                    // 判断key是否过期
	Keys() interface{}                           // 获取所有key
}

// Init 初始化缓存
func Init() {
	// 需要使用的时候直接获取，没有注入到全局变量中
	driver := configDriver()
	register(driver)
}

// GetDefault 获取默认驱动缓存实例
func GetDefault(ctx context.Context) *Cache {
	return GetByDriver(ctx, configDriver())
}

// GetByDriver 获取指定驱动缓存实例: redis memory
func GetByDriver(ctx context.Context, driver string) *Cache {
	driver = strings.ToLower(driver)
	ca, tenant := CtxCache(ctx)

	if ca == nil || ca.driver != driver {
		register(driver)

		ca = &Cache{
			driver: driver,
			store:  Stores[driver],
		}
	}

	ca.SetTenant(ctx, tenant)
	return ca
}

// 实例化一个缓存， 放入缓存池中
func register(driver string) {
	if _, ok := Stores[driver]; ok {
		return
	}

	var store StoreInterface
	if driver == "redis" {
		store = NewRedisStore(config.PlatformAlias)
	} else if driver == "memory" {
		store = NewMemoryStore(config.PlatformAlias)
	} else {
		log.Printf("cache driver \"%s\" not exists. you can change driver is: memory, redis\n", driver)
		return
	}

	Stores[driver] = store
	log.Printf("cache driver \"%s\" connected success.", driver)
}

// 获取配置的缓存驱动
func configDriver() string {
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

func CtxCache(ctx context.Context) (*Cache, string) {
	var ca *Cache
	if caName := ctx.Value("cache"); caName != nil {
		ca = caName.(*Cache)
	}

	tenant := config.PlatformAlias
	if tenantName := ctx.Value("tenant"); tenantName != nil {
		tenant = tenantName.(string)
	}

	return ca, tenant
}

/** 结构体方法 */

func (c *Cache) Driver() string {
	return c.driver
}
func (c *Cache) Tenant() string {
	return c.store.Tenant()
}

func (c *Cache) Platform() bool {
	return c.store.SetTenant(config.PlatformAlias, 0)
}

func (c *Cache) SetTenant(ctx context.Context, tenant string) bool {
	if len(tenant) == 0 {
		_, tenant = CtxCache(ctx)
	}

	// TODO 不可使用
	/*site := model.NewWebSites(ctx).ByUUID("uuid")
	site := 1;*/

	return c.store.SetTenant(tenant, 1)
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
