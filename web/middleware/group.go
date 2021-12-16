package middleware

import (
	"git.kuainiujinke.com/oa/oa-common-golang/web"
)

type group struct {
	Web               web.HandlerFuncGroup
	WebAPI            web.HandlerFuncGroup
	PublicWeb         web.HandlerFuncGroup
	PublicAPI         web.HandlerFuncGroup
	MobileClientAPI   web.HandlerFuncGroup
	OpenAPIClientAuth web.HandlerFuncGroup
	OpenAPIUserAuth   web.HandlerFuncGroup
}

// 中间件组定义
// 底层默认分了几组，若不适用，可调用 AddToGroup 自行定义新的组
var Group = group{
	// 需要 Web 登录态的页面：
	Web: web.HandlerFuncGroup{
		Tenant,
		Cors,
		Permission,
	},
	// 需要 Web 登录态的 AJAX 接口：
	WebAPI: web.HandlerFuncGroup{
		Tenant,
		Cors,
		Permission,
	},
	// 无需登录态的 Web 页面：
	PublicWeb: web.HandlerFuncGroup{
		Tenant,
		Cors,
	},
	// 无需登录态，也无需鉴权的 API 接口：
	PublicAPI: web.HandlerFuncGroup{
		Tenant,
		Cors,
	},
	// 手机 APP 所用的 API，需要【移动客户端登录态/鉴权】
	MobileClientAPI: web.HandlerFuncGroup{
		Tenant,
		APPClientApiAuth,
	},
	// OAuth【用户授权】的 OpenAPI，服务端对服务端
	OpenAPIClientAuth: web.HandlerFuncGroup{
		Tenant,
		OAuthCodeVerify,
	},
	// OAuth【客户端授权】的 OpenAPI，服务端对服务端
	OpenAPIUserAuth: web.HandlerFuncGroup{
		Tenant,
		OAuthClientVerify,
	},
}
