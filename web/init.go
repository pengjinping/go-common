package web

import (
	"net/http"
	"oa-common/config"
	"oa-common/logger"

	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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

	if RouterConfig := config.Get("Routers"); RouterConfig != nil {
		for _, rtc := range RouterConfig.([]config.RouterConfig) {
			routerGroup := Router.Group(rtc.BasePath)
			(*rtc.RouterDefine)(routerGroup)
		}
	}

	return Router
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
