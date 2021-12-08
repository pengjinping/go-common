package config

import "github.com/gin-gonic/gin"

type ConfigType map[string]interface{}

type defaultConfig struct {
	Server ServerConfig `mapstructure:"server" json:"server" yaml:"server"`
	Mysql  MysqlConfig  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	JWT    JWTConfig    `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap    ZapConfig    `mapstructure:"zap" json:"zap" yaml:"zap"`
	Debug  bool         `json:"Debug" yaml:"debug"`
}

type RouterDefine *func(Router *gin.RouterGroup)

type ServerConfig struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type MysqlConfig struct {
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	DBName   string `mapstructure:"dbname" json:"dbname"`
}

type RouterConfig struct {
	BasePath     string
	RouterDefine RouterDefine
}

type ZapConfig struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`                           // 级别
	Format        string `mapstructure:"format" json:"format" yaml:"format"`                        // 输出
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                        // 日志前缀
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`                 // 日志文件夹
	ShowLine      bool   `mapstructure:"show-line" json:"showLine" yaml:"showLine"`                 // 显示行
	EncodeLevel   string `mapstructure:"encode-level" json:"encodeLevel" yaml:"encode-level"`       // 编码级
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktraceKey" yaml:"stacktrace-key"` // 栈名
	LogInConsole  bool   `mapstructure:"log-in-console" json:"logInConsole" yaml:"log-in-console"`  // 输出控制台
}

type JWTConfig struct {
	SigningKey string `mapstructure:"signing-key" json:"signing-key"`
}
