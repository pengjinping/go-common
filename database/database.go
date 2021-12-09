package database

import (
	"context"
	"fmt"
	"log"
	"strings"

	"git.kuainiujinke.com/oa/oa-common/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbConns = make(map[string]*gorm.DB)

func Init() {
	GetDBByName(context.TODO(), "platform")
}

//根据上下文，获取一个合适的 DB 对象
//这里使用原生 Context，而不用 RequestContext，因为不仅仅是 Web 请求才会用到 DB
func GetDB(ctx context.Context) *gorm.DB {
	var tenant string
	if tenantName := ctx.Value("tenant"); tenantName != nil {
		tenant = tenantName.(string)
	} else {
		tenant = "platform"
	}
	return GetDBByName(ctx, tenant)
}

func GetDBByName(ctx context.Context, dbName string) *gorm.DB {
	if _, ok := dbConns[dbName]; !ok {
		conf := getConfig()
		if conf.DBName != dbName && dbName != "" && dbName != "platform" && !strings.Contains(dbName, ":") {
			conf.DBName = dbName
		}
		connect(&conf, dbName)
	}
	return dbConns[dbName].WithContext(ctx)
}

func getConfig() config.MysqlConfig {
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
	if config.GetBool("Debug") {
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

	log.Printf("db \"%s\" connected success", dbName)
}
