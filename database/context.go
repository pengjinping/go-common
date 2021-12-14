package database

import (
	"context"
)

// 返回一个空的【可控的】上下文，内置了一个空的 DB 会话连接池
// 适用于：非 WEB 请求的场景，如 纯异步执行 JOB，在程序开始之前初始化一个上下文，在调用链中传导下去
// tenantUUID: 租户UUID（即为db库名）。当是平台时，固定传常量： config.PlatformAlias
func EmptyContext(tenantUUID string) context.Context {
	c, _ := context.WithCancel(context.WithValue(context.WithValue(context.Background(), CtxPoolKey, NewPool()), "tenant", tenantUUID))
	return c
}
