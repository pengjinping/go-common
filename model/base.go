package model

import (
	"context"

	"git.kuainiujinke.com/oa/oa-common-golang/cache"
	"git.kuainiujinke.com/oa/oa-common-golang/config"
	"git.kuainiujinke.com/oa/oa-common-golang/database"
	"git.kuainiujinke.com/oa/oa-common-golang/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 所有 Model 的公用能力
// (属性名故意写得很长，以减少和 table 业务字段重名的概率)

type BaseModel struct {
	dbConnection   *gorm.DB
	cacheConn      *cache.Cache
	forcePlatform  bool
	currentContext context.Context
}

// 对 model 的一般属性的初始化
func (m *BaseModel) initModel(ctx context.Context) {
	m.currentContext = ctx
}

// 对 model 中所用缓存引擎的初始化（默认 redis 缓存）
func (m *BaseModel) initCache() {
	m.cacheConn = cache.Get(m.currentContext)
}

// 对 model 中使用内存缓存引擎的初始化 （与 initCache() 二选一，不支持两者同时用）
func (m *BaseModel) initMemCache() {
	m.cacheConn = cache.ByDriver(m.currentContext, "memory")
}

// 指定使用平台库
//（不会更改 ctx 中的租户信息，只是本 model 内部的 db 连接变化）

func (m *BaseModel) UsePlatform() {
	m.forcePlatform = true
	m.dbConnection = nil
	if m.cacheConn != nil {
		m.cacheConn.UsePlatform()
	}
}

// 指定使用【传入的】租户库
//（不会更改 ctx 中的租户信息，只是本 model 内部的 db 连接变化）

func (m *BaseModel) UseTenant(tenantUUID string) {
	m.forcePlatform = false
	db := database.ByName(m.currentContext, tenantUUID)
	m.dbConnection = db
	if m.cacheConn != nil {
		m.cacheConn.UseTenant(tenantUUID)
	}
}

// 指定使用【默认的】库 (从 ctx 中推断)
// 可在调用了 UsePlatform()/UseTenant() 之后，调用本方法进行恢复
//（不会更改 ctx 中的租户信息，只是本 model 内部的 db 连接变化）

func (m *BaseModel) UseDefault() {
	m.forcePlatform = false
	m.dbConnection = nil
	if m.cacheConn != nil {
		m.cacheConn.UseDefault()
	}
}

// DB 获取默认的 db 连接
func (m *BaseModel) DB() *gorm.DB {
	if m.currentContext == nil {
		// 这里虽没有 panic，但返回了空，可能会引起调用方的 panic，这件事由调用方自行 recover 掉
		logger.Error(&gin.Context{}, "Model 初始化时，没有上下文信息 (BaseModel.currentContext)")
		return nil
	}
	if m.dbConnection != nil {
		return m.dbConnection
	}
	var db *gorm.DB
	if m.forcePlatform {
		db = database.ByName(m.currentContext, config.PlatformAlias)
	} else {
		db = database.DB(m.currentContext)
	}
	m.dbConnection = db
	return db
}

// Cache 获取cache连接
func (m *BaseModel) Cache() *cache.Cache {
	if m.cacheConn == nil {
		m.cacheConn = cache.Get(m.currentContext)
	}

	return m.cacheConn
}
