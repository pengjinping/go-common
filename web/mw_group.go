package web

import (
	"github.com/gin-gonic/gin"
)

//路由组 所用的 action 函数组
type HandlerFuncGroup []func() gin.HandlerFunc

type middlewareGroup map[string]HandlerFuncGroup

//中间件组定义
//底层默认分了几组，若不适用，可调用 AddToGroup 自行定义新的组
var group = middlewareGroup{
	//需要 Web 登录态的页面：
	"web": HandlerFuncGroup{
		TenantMiddleware,
		CorsMiddleware,
		PermissionMiddleware,
	},
	//需要 Web 登录态的 AJAX 接口：
	"webAPI": HandlerFuncGroup{
		TenantMiddleware,
		CorsMiddleware,
		PermissionMiddleware,
	},
	//无需登录态的 Web 页面：
	"publicWeb": HandlerFuncGroup{
		TenantMiddleware,
		CorsMiddleware,
	},
	//无需登录态，也无需鉴权的 API 接口：
	"publicAPI": HandlerFuncGroup{
		TenantMiddleware,
		CorsMiddleware,
	},
	//手机 APP 所用的 API，需要【移动客户端登录态/鉴权】
	"mobileClientAPI": HandlerFuncGroup{
		TenantMiddleware,
		ApiAuthMiddleware,
	},
	//需要鉴权的 API，服务端对服务端
	"api": HandlerFuncGroup{
		TenantMiddleware,
		ApiAuthMiddleware,
	},
	//需要 OAuth 鉴权的 OpenAPI，服务端对服务端
	"openAPI": HandlerFuncGroup{
		TenantMiddleware,
		ApiAuthMiddleware,
	},
}

//为 路由组 注册 中间件组
func ApplyMiddlewareGroup(e *gin.RouterGroup, groupName string) {
	if funcs, ok := group[groupName]; ok {
		for _, f := range funcs {
			e.Use(f())
		}
	}
}

//向默认的 中间件组 新增内容
func AddToMiddlewareGroup(groupName string, funcGroup HandlerFuncGroup) {
	if _, has := group[groupName]; has {
		group[groupName] = append(group[groupName], funcGroup...)
	} else {
		group[groupName] = funcGroup
	}
}

//获取某个 中间件组 的内容
func GetMiddlewareGroup(groupName string) HandlerFuncGroup {
	if group, has := group[groupName]; has {
		return group
	} else {
		return make(HandlerFuncGroup, 0)
	}
}
