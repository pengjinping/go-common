package web

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.kuainiujinke.com/oa/oa-common/config"
	"git.kuainiujinke.com/oa/oa-common/utils"

	"fmt"

	"github.com/gin-gonic/gin"
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
	//TODO 包住所有运行时panic？

	port := config.GetInt("Server.Port")
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: e,
	}

	go func() {
		log.Printf("启动服务器, 端口： %d", port)

		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("Listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("==================")
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("服务已关闭")
}
