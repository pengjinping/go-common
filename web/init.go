package web

import (
	"net/http"

	"git.kuainiujinke.com/oa/oa-common/config"
	"git.kuainiujinke.com/oa/oa-common/utils"

	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type routerConfig struct {
	MiddlewareGroupName string
	RouterFunc          config.RouterDefine
}

var routerGroups = make(map[string]routerConfig)

func Init() *gin.Engine {
	Router := gin.Default()

	if exists, err := utils.PathExists("./templates"); err == nil && exists {
		Router.LoadHTMLGlob("templates/*")
	}

	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	// 通过 Register注册路由组, 以及各组路由对应的中间件组
	for basePath, config := range routerGroups {
		group := Router.Group(basePath)
		ApplyMiddlewareGroup(group, config.MiddlewareGroupName)
		(*config.RouterFunc)(group)
	}

	return Router
}

// Register 注册路由地址
// @param middlewareGroupname 中间件组的名称，如：web、api、openAPI、publicWeb...等。详见 mw_group.go 中的定义。
// 若没找到合适的中间件组，可以先调用 web.AddToMiddlewareGroup() 来注册自己的组。
func Register(basePath string, middlewareGroupname string, routerFunc func(Router *gin.RouterGroup)) {
	routerGroups[basePath] = routerConfig{middlewareGroupname, &routerFunc}
}

func Start(e *gin.Engine) {
	//TODO 优雅重启/退出
	//TODO 包住所有运行时panic
	port := config.GetInt("Server.Port")
	zap.S().Debugf("启动服务器, 端口： %d", port)
	if err := e.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
}
