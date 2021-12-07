package web

import (
	"context"
	"net/http"
	"oa-common/database"
	"oa-common/model"

	"github.com/gin-gonic/gin"
)

type TenantInfo struct {
	Name string // 数据库连接名称
	UUID string // 租户id
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
	database.GetDB(context.TODO()).Find(&websites)
	tenants := make([]TenantInfo, 0)
	for _, item := range websites {
		tenants = append(tenants, TenantInfo{
			Name: item.Name,
			UUID: item.UUID,
		})
	}
	return tenants
}

func (provider *OaTenantProvider) TenantUUIDResolver(ctx *gin.Context) string {
	return ctx.Request.Host
}

func TenantMiddleware() gin.HandlerFunc {
	// 初始化平台和租户连接
	return func(c *gin.Context) {
		uuid := new(OaTenantProvider).TenantUUIDResolver(c)
		c.Set("tenant", uuid)
		if db := database.GetDB(c); db != nil {
			c.Set("DB", db)
			c.Next()
		} else {
			FailWithMessage(http.StatusNotImplemented, "不存在的租户:"+uuid, c)
			c.Abort()
			return
		}
	}
}