package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type DBPool map[string]*gorm.DB

const CtxPoolKey = "DBPool"

// 初始化一个连接池map
func NewPool() DBPool {
	return make(DBPool)
}

// 获取连接池中某个DB名称所对应的连接key
func Key(dbName string) string {
	return fmt.Sprintf("DB_%s", dbName)
}

// 从上下文中尝试获取一个 db 连接
func FromCtx(c context.Context, dbName string) *gorm.DB {
	if p := poolFromCtx(c); p == nil {
		return nil
	} else {
		return FromPool(p, dbName)
	}
}

// 从给定连接池中尝试获取一个 db 连接
func FromPool(p DBPool, dbName string) *gorm.DB {
	if db := p[Key(dbName)]; db != nil {
		return db
	}
	return nil
}

// 向上下文的连接池中放入一个 db 连接
func SetCtxDB(c context.Context, dbName string, conn *gorm.DB) {
	if p := poolFromCtx(c); p != nil {
		SetPoolDB(p, dbName, conn)
	}
}

// 向指定连接池中放入一个 db 连接
func SetPoolDB(p DBPool, dbName string, conn *gorm.DB) {
	p[Key(dbName)] = conn
}

func poolFromCtx(c context.Context) DBPool {
	var p DBPool
	pool := c.Value(CtxPoolKey)
	if pool == nil {
		return nil
	}
	p = pool.(DBPool)
	return p
}
