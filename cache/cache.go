package cache

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"git.kuainiujinke.com/oa/oa-common-golang/config"
	"git.kuainiujinke.com/oa/oa-common-golang/web"
)

/**
// 获取本次请求的默认缓存 【支持多租户分组】
ca := cache.Get(c)
ca := cache.GetByDriver(c, "memory")	// 也可以指定驱动 如: 指定缓存驱动示例

ca.Driver()		// 获取缓存驱动名称
ca.Tenant() 	// 获取缓存租户uuid

// 切换租户信息
ca.UsePlatform()     // 切换到平台
ca.UseDefault()   	 // 切回当前租户
ca.UseTenant(uuid)   // 切换到指定租户

// 设置缓存，有有效期 单位是s  当时间为0时 是永久有效 等于Forever
ca.Set("AA", "ABC", 2)
ca.Forever("AA", "SSSS")	// 设置永久缓存

ca.Get("AA")		// 获取缓存
ca.Has("AA")		// 是否存在
ca.IsExpire("AA")	// 是否过期
ca.Delete("AA")		// 删除缓存
ca.Keys()			// 获取所有缓存keys

// 缓存支持闭包函数 在闭包中可以通过后面的参数传进去 不要使用全局变量防止数据污染
res := ca.Remember("key", 5, func(args ...interface{}) (interface{}, error) {
	a := args[0].(int)
	b := args[1].(int)
	return a + b + 456, nil
}, 4, 500)
*/

var Stores = make(map[string]StoreInterface)
var lock sync.Mutex

type Cache struct {
	driver         string // 缓存驱动
	currentContext context.Context
	store          StoreInterface
}

type StoreInterface interface {
	Tenant() string                              // 获取租户UUID
	SetTenant(tenant string, tenantId int) bool  // 设置连接库DB
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

// Get 获取默认驱动缓存实例
func Get(ctx context.Context) *Cache {
	return ByDriver(ctx, configDriver())
}

// ByDriver 获取指定驱动缓存实例: redis memory
func ByDriver(ctx context.Context, driver string) *Cache {
	driver = strings.ToLower(driver)

	var ca *Cache
	if caName := ctx.Value("cache"); caName != nil {
		ca = caName.(*Cache)
	}

	if ca == nil || ca.driver != driver {
		register(driver)

		ca = &Cache{
			driver:         driver,
			currentContext: ctx,
			store:          Stores[driver],
		}
	}

	ca.UseDefault()
	return ca
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

// 实例化一个缓存， 放入缓存池中
func register(driver string) {
	if _, ok := Stores[driver]; ok {
		return
	}

	lock.Lock()

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
	lock.Unlock()

	log.Printf("cache driver \"%s\" connected success.", driver)
}

/** 结构体方法 */

func (c *Cache) Driver() string {
	return c.driver
}
func (c *Cache) Tenant() string {
	return c.store.Tenant()
}

// 切换租户信息

func (c *Cache) UsePlatform() bool {
	return c.UseTenant(config.PlatformAlias)
}

func (c *Cache) UseDefault() bool {
	tenant := config.PlatformAlias
	if tenantName := c.currentContext.Value("tenant"); tenantName != nil {
		tenant = tenantName.(string)
	}

	return c.UseTenant(tenant)
}

func (c *Cache) UseTenant(tenant string) bool {
	if len(tenant) == 0 {
		log.Printf("缓存：切换租户失败，租户UUID不可为空")
		return false
	}

	site := web.SiteByUUID(tenant)
	if site == nil {
		log.Printf("缓存：切换租户失败，租户%s不存在", tenant)
		return false
	}

	return c.store.SetTenant(tenant, int(site.ID))
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
	cacheValue := c.store.Get(key)
	if cacheValue != nil {
		return cacheValue
	}

	value, err := f(args...)
	if err != nil {
		return nil
	}

	c.Set(key, value, time)

	return value
}
