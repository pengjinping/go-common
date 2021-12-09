package web

import (
	"net/http"

	"git.kuainiujinke.com/oa/oa-common/config"
	"git.kuainiujinke.com/oa/oa-common/logger"

	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var routerGroups = make(map[string]config.RouterDefine)

func Init() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	//TODO 中间件组
	Router.Use(TenantMiddleware())
	Router.Use(RecoveryMiddleware(logger.Logger, true))

	// 通过 Register注册路由
	for basePath, routeFunc := range routerGroups {
		(*routeFunc)(Router.Group(basePath))
	}

	return Router
}

// Register 注册路由地址
func Register(basePath string, routeFunc func(Router *gin.RouterGroup)) {
	routerGroups[basePath] = &routeFunc
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
