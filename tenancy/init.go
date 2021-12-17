package tenancy

/*
之所以本包单独成包，而不放在 model 包中，是为了规避循环依赖
因为 model 包依赖 cache 包，而 cache 包又依赖本包
*/

import (
	"git.kuainiujinke.com/oa/oa-common-golang/config"
	"git.kuainiujinke.com/oa/oa-common-golang/database"
)

// 设置所有租户信息，缓存在内存中
// TODO 这里使用内存缓存，没有过期时间机制
// 		当新增租户的时候，需要主动再调用这个方法，来刷新这里的缓存

func Init() {
	var sites []Websites
	Tenants = make(map[string]Websites)

	database.DB(database.EmptyContext(config.PlatformAlias)).Find(&sites)
	for _, v := range sites {
		Tenants[v.UUID] = v
	}

	Tenants[config.PlatformAlias] = Websites{
		ID:   0,
		Name: config.PlatformAlias,
		UUID: config.GetString("server.host"),
	}

}
