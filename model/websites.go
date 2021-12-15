package model

import (
	"context"
	"fmt"

	"git.kuainiujinke.com/oa/oa-common-golang/database"
	"gorm.io/gorm"
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
// ====================================
// NewWebSites() 获取一个 租户表 Model实体
func NewWebSites(ctx context.Context) *Websites {
	var w Websites
	w.forcePlatform = true // 强制使用系统库。若无需强制，请删除此行
	w.currentContext = ctx
	return &w
}

// 注意这里接收 w Websites 而非 w *Websites
// 好处是，会将当前 model 的 BaseModel 基础属性（如 db 连接） 传导到返回值中
// ====================================
// ByUUID() 根据 UUID 获取对应一行结果
func (w Websites) ByUUID(uuid string) *Websites {
	w.DB().Where("uuid = ?", uuid).First(&w)

	return &w
}

// 事务用法的示例
func (w *Websites) TransactionDemo() {
	db := w.DB()
	var funcs []database.TxFunc
	funcs = append(funcs, func(db *gorm.DB) error {
		n := database.Selected(db)
		fmt.Printf("\n===============事务方法-1，DB名称：%s\n", n)
		return nil
	})
	funcs = append(funcs, func(db *gorm.DB) error {
		n := database.Selected(db)
		fmt.Printf("\n===============事务方法-2，DB名称：%s\n", n)
		return nil
	})
	if err := database.Tx(db, funcs...); err != nil {
		fmt.Println(err.Error())
	}
}
