package initialize

import (
	"git.kuainiujinke.com/oa/oa-common/cache"
	"git.kuainiujinke.com/oa/oa-common/config"
	"git.kuainiujinke.com/oa/oa-common/database"
	"git.kuainiujinke.com/oa/oa-common/logger"
	"git.kuainiujinke.com/oa/oa-common/web"

	"github.com/gin-gonic/gin"
)

func InitWebEngine(c *config.ConfigType) *gin.Engine {
	Init(c)
	//初始化路由
	return web.Init()
}

func Init(c *config.ConfigType) {
	//0. 初始化配置
	config.Init(c)

	//1. 初始化logger
	logger.Init()

	//2. 初始化数据库信息
	database.Init()

	//3. 初始化缓存
	cache.Init()
}
