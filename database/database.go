package database

import (
	"context"
	"fmt"
	"log"
	"oa-common/config"
	"strings"

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
	fmt.Println(ctx.Value("tenant"))
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
		conf := config.Get("Mysql").(config.MysqlConfig)
		if conf.DBName != dbName && dbName != "" && dbName != "platform" && !strings.Contains(dbName, ":") {
			conf.DBName = dbName
		}
		connect(&conf, dbName)
	}
	return dbConns[dbName].WithContext(ctx)
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

// func InitDB() {
// 	mysqlConfig := core_global.EntryData.MysqlConfig
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 		mysqlConfig.Username, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.PlatformDbName)
// 	newLogger := logger.New(
// 		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
// 		logger.Config{
// 			SlowThreshold: time.Second, // 慢 SQL 阈值
// 			LogLevel:      logger.Info, // Log level
// 			Colorful:      true,        // 禁用彩色打印
// 		},
// 	)
// 	var err error
// 	core_global.PlatformDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
// 		NamingStrategy: schema.NamingStrategy{
// 			SingularTable: true,
// 		},
// 		Logger: newLogger,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func test() {
// 	type Product struct {
// 		Model
// 		Code  string
// 		Price uint
// 	}
// 	db := GetDB(context.TODO())

// 	// 迁移 schema
// 	db.AutoMigrate(&Product{})

// 	// Create
// 	db.Create(&Product{Code: "D42", Price: 100})

// 	// Read
// 	var product Product
// 	db.First(&product, 1)                 // 根据整形主键查找
// 	db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
// 	fmt.Printf("%v", product)

// 	// Update - 将 product 的 price 更新为 200
// 	db.Model(&product).Update("Price", 200)
// 	// Update - 更新多个字段
// 	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
// 	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

// 	// Delete - 删除 products
// 	db.Delete(&product, 1)
// }
