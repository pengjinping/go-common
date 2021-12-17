package model

import (
	"context"
	"fmt"
	"strings"

	"git.kuainiujinke.com/oa/oa-common-golang/config"
	"git.kuainiujinke.com/oa/oa-common-golang/database"
	"gorm.io/gorm"
)

type KeyValue struct {
	BaseModel
	Value string `json:"value" gorm:"comment:value"`
}

func (KeyValue) TableName() string {
	return "keyvalue"
}

// Model 的用法
// 业务方
// 注意：
// 每个 Model 都【必须】有一个对应的 New 方法
// 1、方便这样的链式调用： model.NewKeyValue(ctx).KeyValue("uuid")
// 2、若是“默认强制使用平台库”，这里可以预先通过 UsePlatform() 进行设置

func NewKeyValue(ctx context.Context) *KeyValue {
	var kv KeyValue
	kv.InitModel(ctx)
	kv.InitCache()
	//kv.UsePlatform() // 强制使用系统库。若无需强制，请删除此行
	return &kv
}

func (m *KeyValue) KeyValue(key string) string {
	ca := m.Cache()
	if !ca.IsExpire(key) {
		return ca.Get(key).(string)
	}

	arr := strings.Split(key, ".")
	condition := map[string]interface{}{"status": "active"}

	if len(arr) == 2 {
		condition["key"] = arr[1]
		condition["app"] = arr[0]
	} else if len(arr) == 1 {
		condition["key"] = key
	} else {
		return ""
	}

	m.DB().Where(condition).First(&m)
	ca.Set(key, m.Value, config.KeyValueDefaultExpire)

	return m.Value
}

// 事务用法的示例
func (m *KeyValue) TransactionDemo() {
	db := m.DB()
	var funcs []database.TxFunc
	funcs = append(funcs, func(tx *gorm.DB) error {
		n := database.Selected(tx)
		fmt.Printf("\n===============事务方法-1，DB名称：%s\n", n)
		return nil
	})
	funcs = append(funcs, func(tx *gorm.DB) error {
		n := database.Selected(tx)
		fmt.Printf("\n===============事务方法-2，DB名称：%s\n", n)
		return nil
	})
	if err := database.Tx(db, funcs...); err != nil {
		fmt.Println(err.Error())
	}
}
