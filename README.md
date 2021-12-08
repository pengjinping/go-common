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
	
	// 注册路由地址 [应用项目的路由包]
	router.Init()
	
	//...可调用 customConfigs.Set() 方法，设置自己的一些配置项（同名key将会覆盖默认配置）

	e := initialize.InitWebEngine(&customConfigs)

	web.Start(e)
}
```

### 路由注册
router/init.go
```golang
package router
import "oa-common/web"

func Init()  {
	web.Register("/oa-micro-message/v1", InitApiRouter)
}
```
router/api.go
```golang
package router

import (
	"github.com/gin-gonic/gin"
)

func InitApiRouter(Router *gin.RouterGroup) {
	...
}
```