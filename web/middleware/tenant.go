package middleware

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"git.kuainiujinke.com/oa/oa-common-golang/cache"
	"git.kuainiujinke.com/oa/oa-common-golang/config"
	"git.kuainiujinke.com/oa/oa-common-golang/database"
	"git.kuainiujinke.com/oa/oa-common-golang/logger"
	"git.kuainiujinke.com/oa/oa-common-golang/web"

	"github.com/gin-gonic/gin"
)

func TenantUUIDResolver(ctx *gin.Context) string {
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

func Tenant() gin.HandlerFunc {
	// 初始化平台和租户连接
	return func(c *gin.Context) {
		// 初始化【本请求专用的】db连接池
		c.Set(database.CtxPoolKey, make(database.DBPool))

		uuid := TenantUUIDResolver(c)
		c.Set("tenant", uuid)

		// 切换缓存
		if ca := cache.Get(c); ca != nil {
			c.Set("cache", ca)
		}

		// 切换日志
		if log := logger.Get(c); log != nil {
			c.Set("logger", log)
			c.Set("loggerId", fmt.Sprintf("%v-%v", time.Now().UnixMicro(), rand.Intn(9999)))
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
