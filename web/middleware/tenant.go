package middleware

import (
	"context"
	"strings"

	"git.kuainiujinke.com/oa/oa-common-golang/cache"
	"git.kuainiujinke.com/oa/oa-common-golang/config"
	"git.kuainiujinke.com/oa/oa-common-golang/database"
	"git.kuainiujinke.com/oa/oa-common-golang/model"
	"git.kuainiujinke.com/oa/oa-common-golang/web"

	"github.com/gin-gonic/gin"
)

type TenantInfo struct {
	ID   uint   //租户ID
	Name string // 数据库连接名称/域名
	UUID string // 租户UUID
}

type TenantProvider interface {
	// 租户的数据提供者
	TenantsProvider() []TenantInfo
	// 租户UUID解析器
	TenantUUIDResolver(*gin.Context) string
}

type OaTenantProvider struct{}

func (provider *OaTenantProvider) TenantsProvider() []TenantInfo {
	var websites []model.Websites
	database.DB(context.TODO()).Find(&websites)
	tenants := make([]TenantInfo, 0)
	for _, item := range websites {
		tenants = append(tenants, TenantInfo{
			ID:   item.ID,
			Name: item.Name,
			UUID: item.UUID,
		})
	}
	return tenants
}

func (provider *OaTenantProvider) TenantUUIDResolver(ctx *gin.Context) string {
	h := ctx.Request.Host
	// 若有端口号，去除之
	if i := strings.Index(h, ":"); i > -1 {
		h = h[:i]
	}
	// 若是对平台的请求，则 uuid 设为固定的别名
	if h == config.GetString("server.host") {
		h = config.PlatformAlias
	}
	return h
}

// 设置所有租户信息  [UUID]ID
func (provider *OaTenantProvider) ConfigWebsite() {
	if _, ok := config.WebSite[config.PlatformAlias]; !ok {
		config.WebSite[config.PlatformAlias] = 0
		tenants := provider.TenantsProvider()
		for _, item := range tenants {
			config.WebSite[item.UUID] = int(item.ID)
		}
	}
}

func Tenant() gin.HandlerFunc {
	// 初始化平台和租户连接
	return func(c *gin.Context) {
		new(OaTenantProvider).ConfigWebsite()

		// 初始化【本请求专用的】db连接池
		c.Set(database.CtxPoolKey, make(database.DBPool))

		uuid := new(OaTenantProvider).TenantUUIDResolver(c)
		c.Set("tenant", uuid)
		// TODO set 租户的 ID

		// 切换缓存
		if ca := cache.GetDefault(c); ca != nil {
			c.Set("cache", ca)
		}

		// 切换数据库
		if db := database.DB(c); db != nil {
			c.Next()
		} else {
			web.FailWithMessage("不存在的租户:"+uuid, c)
			c.Abort()
			return
		}
	}

}
