package model

import (
	"context"
	"git.kuainiujinke.com/oa/oa-common-golang/cache"
	"git.kuainiujinke.com/oa/oa-common-golang/config"
	"strings"
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
// 2、若是“默认强制使用平台库”，这里会预设 forcePlatform 为 true

func NewKeyValue(ctx context.Context) *KeyValue {
	var kv KeyValue
	kv.currentContext = ctx
	kv.cacheConn = cache.Get(ctx)
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
