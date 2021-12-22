# go-common

GO公共包，不独立运行，而是被 import 到业务应用中。

提供能力：
- 配置
- 多租户[SAAS]
  - 日志
  - DB
  - CACHE
- KV
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
	"github.com/pengjinping/go-common/config"
	"github.com/pengjinping/go-common/initialize"
	"github.com/pengjinping/go-common/web"
)

func main() {
	var customConfigs config.ConfigType
	
	// 注册路由地址 [应用项目的路由包]
	router.Init()
	
	e := initialize.InitWebEngine(&customConfigs)

	web.Start(e)
}
```

### 路由注册
router/init.go
```golang
package router
import "git.kuainiujinke.com/oa/oa-common-golang/web"

func Init()  {
  web.Register("/message/open_api", middleware.Group.PublicAPI, InitOpenAPIRouter)
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

### 缓存使用
在配置文件中设置缓存驱动,目前支持：memory(内存缓存), redis
```yaml
cache:
  driver: memory
```
