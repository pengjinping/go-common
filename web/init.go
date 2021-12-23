package web

import (
	"context"
	"errors"
	"git.kuainiujinke.com/oa/oa-common-golang/utils/directory"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fmt"
	"git.kuainiujinke.com/oa/oa-common-golang/config"

	"github.com/gin-gonic/gin"
)

// 路由组的函数类型
type RouterDefine func(Router *gin.RouterGroup)

// 路由组 所用的 action 函数组
type HandlerFuncGroup []func() gin.HandlerFunc

// 路由组的预定义结构
type routerConfig struct {
	MiddlewareGroup HandlerFuncGroup
	RouterFunc      RouterDefine
}

var routerGroups = make(map[string]routerConfig)

func Init() *gin.Engine {
	var mode string
	switch config.GetString("env") {
	case "prod":
		mode = gin.ReleaseMode
	case "test":
		mode = gin.TestMode
	default:
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	e := gin.Default()

	if exists, err := directory.PathExists("./templates"); err == nil && exists {
		e.LoadHTMLGlob("templates/*")
	}

	e.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	// 向 gin 引擎中注册预设的路由组, 以及各组路由对应的中间件组
	for basePath, config := range routerGroups {
		group := e.Group(basePath)
		ApplyMiddleware(group, config.MiddlewareGroup)
		(config.RouterFunc)(group)
	}

	return e
}

// 为 路由组 注册 中间件组
func ApplyMiddleware(e *gin.RouterGroup, g HandlerFuncGroup) {
	for _, f := range g {
		e.Use(f())
	}
}

// Register 注册路由地址
// @param middlewareGroup 中间件组，默认提供的如：Web、API、OpenAPI、PublicWeb...等。详见 middleware.Group 中的定义。
// 若没找到合适的中间件组，可以传入自己的组。
func Register(basePath string, middlewareGroup HandlerFuncGroup, routerFunc RouterDefine) {
	routerGroups[basePath] = routerConfig{middlewareGroup, routerFunc}
}

func Start(e *gin.Engine) {
	//TODO 包住所有运行时panic？

	port := config.GetInt("Server.Port")
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: e,
	}

	go func() {
		log.Printf("启动服务器, 端口：%d", port)

		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("Listen %s\n", err)
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
