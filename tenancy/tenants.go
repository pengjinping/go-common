package tenancy

/*
之所以本包单独成包，而不放在 model 包中，是为了规避循环依赖
因为 model 包依赖 cache 包，而 cache 包又依赖本包
*/

type Websites struct {
	ID   uint   `gorm:"primarykey"`
	Name string `json:"name" gorm:"comment:name"`
	UUID string `json:"uuid" gorm:"comment:uuid"`
}

// 用于存储各个租户列表信息
var Tenants map[string]Websites

// ByUUID() 根据 UUID 获取对应一行结果
func ByUUID(UUID string) *Websites {
	if w, ok := Tenants[UUID]; ok {
		return &w
	}

	return nil
}

// ByID() 根据 ID 获取对应一行结果
func ByID(ID uint) *Websites {
	for _, v := range Tenants {
		if v.ID == ID {
			return &v
		}
	}

	return nil
}
