package model

import (
	"context"
)

type Websites struct {
	BaseModel
	ID   uint   `gorm:"primarykey"`
	Name string `json:"name" gorm:"comment:name"`
	UUID string `json:"uuid" gorm:"comment:uuid"`
}

// Model 的用法
// 业务方
// 注意：
// 每个 Model 都【必须】有一个对应的 New 方法
// 1、方便这样的链式调用： model.NewWebSites(ctx).ByUUID("uuid")
// 2、若是“默认强制使用平台库”，这里会预设 forcePlatform 为 true

// 获取一个 租户表 Model实体
func NewWebSites(ctx context.Context) *Websites {
	var w Websites
	w.forcePlatform = true // 强制使用系统库。若无需强制，请删除此行
	w.currentContext = ctx
	return &w
}

func (w Websites) ByUUID(uuid string) *Websites {
	w.DB().Where("uuid = ?", uuid).First(&w)

	return &w
}
