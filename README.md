# oa-common-golang

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
	"git.kuainiujinke.com/oa/oa-common-golang/config"
	"git.kuainiujinke.com/oa/oa-common-golang/initialize"
	"git.kuainiujinke.com/oa/oa-common-golang/web"
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
import "git.kuainiujinke.com/oa/oa-common-golang/web"

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

### 缓存使用
在配置文件中设置缓存驱动,目前支持：memory(内存缓存), redis
```yaml
cache:
  driver: memory
```
缓存的使用：
```golang
    // 获取本次请求的默认缓存 【支持多租户分组】
    cache := cache.GetDefault(c)
    
    // 也可以指定驱动 如: 指定缓存驱动示例
    caMemory := cache.GetByDriver(c, "memory")
    
    // 获取缓存驱动名称
    fmt.Printf("%s\n", ca.GetStoreName())
    
    // 获取缓存租户uuid
    fmt.Printf("%s\n", ca.GetTenant())

    // 设置缓存，有有效期 单位是s  当时间为0时 是永久有效 等于Forever
    ca.Set("AA", timehelper.FormatDateTime(time.Now()), 2)
    // 设置永久缓存 
    ca.Forever("AA", timehelper.FormatDateTime(time.Now()))
    
    // 获取缓存、是否存在、是否过期
    fmt.Printf("测试值：%s\n", ca.Get("AA"))
    fmt.Printf("是否存在：%v\n", ca.Has("AA"))
    fmt.Printf("是否过期：%v\n", ca.IsExpire("AA"))
    
    // 删除缓存
    ca.Delete("AA")
    
    // 缓存支持闭包函数 在闭包中可以通过后面的参数传进去 不要使用全局变量防止数据污染
    res := ca.Remember("key", 5, func(args ...interface{}) (interface{}, error) {
        a := args[0].(int)
        b := args[1].(int)
        return a + b + 456, nil
    }, 4, 500)
```