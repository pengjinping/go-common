package model

import (
	"context"

	"git.kuainiujinke.com/oa/oa-go-common/config"
	"git.kuainiujinke.com/oa/oa-go-common/database"
	"git.kuainiujinke.com/oa/oa-go-common/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 所有 Model 的公用能力
// (属性名故意写得很长，以减少和 table 业务字段重名的概率)
type BaseModel struct {
	dbConnection   *gorm.DB
	forcePlatform  bool
	currentContext context.Context
}

// 获取默认的 db 连接
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
