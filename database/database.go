package database

import (
	"context"
	"fmt"
	"log"
	"strings"

	"git.kuainiujinke.com/oa/oa-common-golang/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbConns = NewPool()

func Init() {
	ByName(context.TODO(), config.PlatformAlias)
}

//根据上下文，获取一个合适的 DB 连接会话
//这里使用原生 Context，而不用 RequestContext，因为不仅仅是 Web 请求才会用到 DB
func DB(ctx context.Context) *gorm.DB {
	var tenant string
	if tenantName := ctx.Value("tenant"); tenantName != nil {
		tenant = tenantName.(string)
	} else {
		tenant = config.PlatformAlias
	}
	return ByName(ctx, tenant)
}

// 根据 DB名称，获取一个 db 连接会话
// 若上下文连接池中已有，则会复用之
// 若是平台库，dbName 固定传常量：config.PlatformAlias
func ByName(ctx context.Context, dbName string) *gorm.DB {
	if db := FromCtx(ctx, dbName); db != nil {
		return db
	}

	if _, ok := dbConns[dbName]; !ok {
		conf := Config()
		if conf.DBName != dbName && dbName != "" && dbName != config.PlatformAlias && !strings.Contains(dbName, ":") {
			conf.DBName = dbName
		}
		connect(&conf, dbName)
	}

	newDB := dbConns[dbName].WithContext(ctx)
	SetCtxDB(ctx, dbName, newDB)
	return newDB
}

// 获取一个连接所 select 的库名称 （Mysql）
func Selected(db *gorm.DB) string {
	var name string
	db.Raw("select database()").Scan(&name)
	return name
}

func Config() config.MysqlConfig {
	var conf config.MysqlConfig
	if err := config.UnmarshalKey("Mysql", &conf); err != nil {
		log.Printf("DB Config init failed: %s\n", err)
	}
	return conf
}

func connect(cfg *config.MysqlConfig, dbName string) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	var ormLogger logger.Interface
	// 生产环境 强制不打印详细 SQL 日志
	if config.GetBool("debug") && !config.IsProduction() {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: ormLogger,
	})
	if err != nil {
		log.Fatal(err)
	}

	dbConns[dbName] = db

	log.Printf("DB \"%s\" connected success", dbName)
}
