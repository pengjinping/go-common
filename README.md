# OA-common

OA 公共包，不独立运行，而是被 import 到业务应用中。

提供能力：
- 日志
- 多租户
- DB
- 配置
- Web 引擎 （基于 gin 框架）
- ...

## 使用
### 配置文件
开发环境配置，详见：`deployment/application-dev.yaml`

//TODO 完善配置用法说明（ENV 等）

### main.go
```golang
package main

import (
	"oa-common/config"
	"oa-common/initialize"
	"oa-common/web"
)

func main() {
	var customConfigs config.ConfigType

	//应用的个性化路由组，开始
	var routerConfig = make([]config.RouterConfig, 0)
	apiRouter := router.InitApiRouter
	routerConfig = append(routerConfig, config.RouterConfig{
		BasePath:      "/oa-micro-message/v1",
		RouterDefine: &apiRouter,
	})
	//...
	customConfigs.Set("Routers", routerConfig)
	//应用的个性化路由组，结束

	//...可调用 customConfigs.Set() 方法，设置自己的一些配置项（同名key将会覆盖默认配置）

	e := initialize.InitWebEngine(&customConfigs)

	web.Start(e)
}
```