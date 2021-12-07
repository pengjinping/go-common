package web

import (
	"github.com/gin-gonic/gin"
)

type HandlerFuncGroup []func() gin.HandlerFunc
type middlewareGroup map[string]HandlerFuncGroup

//底层默认分了几组，若不适用，可调用 AddToGroup 自行定义新的组
var group = middlewareGroup{
	"api": HandlerFuncGroup{
		ApiAuthMiddleware,
		TenantMiddleware,
	},
	"web": HandlerFuncGroup{
		TenantMiddleware,
		CorsMiddleware,
		PermissionMiddleware,
	},
}

func applyMiddlewareGroup(e *gin.Engine, groupName string) {
	if funcs, ok := group[groupName]; ok {
		for _, f := range funcs {
			e.Use(f())
		}
	}
}

func AddToMiddlewareGroup(groupName string, funcGroup HandlerFuncGroup) {
	if _, has := group[groupName]; has {
		group[groupName] = append(group[groupName], funcGroup...)
	} else {
		group[groupName] = funcGroup
	}
}
