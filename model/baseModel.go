package model

import (
	"context"

	"gorm.io/gorm"
)

type BaseModel struct {
}

func (m *BaseModel) getConn(ctx context.Context) *gorm.DB {
	if db := ctx.Value("db"); db.(bool) {
		return db.(*gorm.DB)
	} else {
		panic("获取 DB 连接失败")
	}

}
