package keyvalue

import (
	"context"
	"encoding/json"
	"git.kuainiujinke.com/oa/oa-common-golang/model"
)

type KeyValue struct {
	model *model.KeyValue
	value string
}

/**
 * 业务方调用说明：
 * 如果只是快速获取一个值 keyvalue.NewKeyValue(ctx).Values(key)
 * 	// 快速获取某个key值
 *	kv0 := keyvalue.NewKeyValue(c).Key("employee.employee_image_update_on_off").Bool()
 *
 *	// 如果需要使用多个缓存值 则可以先设置后获取
 *	kv := keyvalue.NewKeyValue(c)
 *	kv1 := kv.Key("employee.employee_image_update_on_off").Bool()
 *
 *	kv2 := kv.Key("employee.roster_download_except").Map()
 *
 *	// 切换获取平台和当前租户 Platform()  Tenant(c)
 *	kv := keyvalue.NewKeyValue(c).Platform().Key("workbenchApps").Bool()
 *  或者
 *  kv := keyvalue.NewKeyValue(c)
 *	kv.Platform()
 *	kv3 := kv.Key("workbenchApps").Bool()
 *
 *	// 切回到当前租户KV
 *  kv.Platform(c)
 */

func NewKeyValue(ctx context.Context) *KeyValue {
	var kv KeyValue
	kv.model = model.NewKeyValue(ctx)
	return &kv
}

func (kv *KeyValue) Tenant() *KeyValue {
	kv.model.UseDefault()
	return kv
}

func (kv *KeyValue) Platform() *KeyValue {
	kv.model.UsePlatform()
	return kv
}

// ---- 获取值信息 -----

func (kv *KeyValue) Key(key string) *KeyValue {
	kv.value = kv.model.KeyValue(key)
	return kv
}

func (kv *KeyValue) Bool() bool {
	return kv.value == "1" || kv.value == "true"
}

func (kv *KeyValue) Values() interface{} {
	return kv.value
}

func (kv *KeyValue) Map() map[string]interface{} {
	var tempMap map[string]interface{}

	err := json.Unmarshal([]byte(kv.value), &tempMap)
	if err != nil {
		panic(err)
	}

	return tempMap
}

func (kv *KeyValue) List() []map[string]interface{} {
	var tempMap []map[string]interface{}

	err := json.Unmarshal([]byte(kv.value), &tempMap)
	if err != nil {
		panic(err)
	}

	return tempMap
}

func (kv *KeyValue) Unmarshal(v interface{}) {
	err := json.Unmarshal([]byte(kv.value), &v)
	if err != nil {
		panic(err)
	}
}
