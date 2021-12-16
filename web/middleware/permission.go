package middleware

import (
	"fmt"
	"git.kuainiujinke.com/oa/oa-common-golang/logger"
	"git.kuainiujinke.com/oa/oa-common-golang/model"
	"git.kuainiujinke.com/oa/oa-common-golang/web"

	"github.com/gin-gonic/gin"
)

func Permission() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取员工的email
		emailMix, e1 := c.Get("email")
		userIdMix, e2 := c.Get("userId")
		if !e1 || !e2 {
			logger.Error(c, "授权用户权限:用户不存在")
			web.FailWithMessage("登录已失效", c)
			c.Abort()
			return
		}

		permissionMap := model.NewPermission(c).UserPermission(emailMix.(string), userIdMix.(string), 600)
		if permissionMap == nil || len(permissionMap) <= 0 {
			logger.Error(c, fmt.Sprintf("获取用户[%v]权限失败", emailMix.(string)))
			web.FailWithMessage("你没有授权任何权限哦", c)
			c.Abort()
			return
		}

		url := c.Request.URL.String()
		if _, ok := permissionMap[url]; ok {
			c.Next()
		} else {
			web.FailWithMessage("禁止访问", c)
			c.Abort()
			return
		}
	}
}
