package web

/*
* 本文件原本是放在 model 包中
* 因为 cache 包中引用了本文件的方法，为避免包循环引用，将本文件从 model 包中移至此处
 */

import (
	"fmt"

	"git.kuainiujinke.com/oa/oa-common-golang/config"
	"git.kuainiujinke.com/oa/oa-common-golang/database"
)

type Websites struct {
	ID   uint   `gorm:"primarykey"`
	Name string `json:"name" gorm:"comment:name"`
	UUID string `json:"uuid" gorm:"comment:uuid"`
}

// 缓存各个租户列表信息
var websitesCache map[string]Websites
var initWebsitesLog = make(map[string]int)

type oneSiteFunc func() *Websites

// 根据 UUID 获取一个租户
func SiteByUUID(UUID string) *Websites {
	f := func() oneSiteFunc {
		return func() *Websites {
			if w, ok := websitesCache[UUID]; ok {
				return &w
			}

			return nil
		}
	}()
	return oneSiteCacheFor(fmt.Sprintf("UUID_%s", UUID), f)
}

// 根据 ID 获取一个租户
func SiteByID(ID uint) *Websites {
	f := func() oneSiteFunc {
		return func() *Websites {
			for _, v := range websitesCache {
				if v.ID == ID {
					return &v
				}
			}

			return nil
		}
	}()
	return oneSiteCacheFor(fmt.Sprintf("ID_%d", ID), f)
}

// 从平台库中查出所有租户信息，缓存在内存中
// TODO 当新增租户的时候，最好是能主动再调用这个方法，来刷新这里的缓存
func initWebsites() {
	var sites []Websites
	websitesCache = make(map[string]Websites)

	database.DB(database.EmptyContext(config.PlatformAlias)).Find(&sites)
	for _, v := range sites {
		websitesCache[v.UUID] = v
	}
	websitesCache[config.PlatformAlias] = Websites{
		ID:   0,
		Name: config.PlatformAlias,
		UUID: config.GetString("server.host"),
	}
	fmt.Println("======INIT SITES=====")
}

// 执行一个从缓存中获取租户的函数，若没获取到，则重新初始化缓存，再尝试一次
func oneSiteCacheFor(reason string, f oneSiteFunc) *Websites {
	res := f()
	if res == nil {
		initWebsitesFor(reason)
		res = f()
	}
	return res
}

// 为某原因执行 InitWebsites，但会防止过多次的执行
func initWebsitesFor(reason string) {
	maxTimes := 3
	if times, ok := initWebsitesLog[reason]; !ok || times < maxTimes {
		initWebsitesLog[reason] = times + 1
		initWebsites()
	}
}
