package database

import (
	"context"
)

// 返回一个空的【可控的】上下文，内置了一个空的 DB 会话连接池
// tenantUUID: 租户UUID（即为db库名）。当是平台时，固定传常量： config.PlatformAlias
func EmptyContext(tenantUUID string) context.Context {
	c, _ := context.WithCancel(context.WithValue(context.WithValue(context.Background(), CtxPoolKey, NewPool()), "tenant", tenantUUID))
	return c
}
