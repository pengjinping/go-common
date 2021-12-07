package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
)

const (
	ConfigEnv      = "GIN_MODE"
	ConfigDir      = "./deployment/"
	DevConfigFile  = "application-dev.yaml"
	TestConfigFile = "application-test.yaml"
	ProdConfigFile = "application-prod.yaml"
)

var Config = make(ConfigType)

func (c *ConfigType) Get(path string) interface{} {
	return funk.Get(&Config, path)
}

func Get(path string) interface{} {
	return Config.Get(path)
}

func GetString(path string) string {
	return Config.Get(path).(string)
}

func GetInt(path string) int {
	return Config.Get(path).(int)
}

func GetBool(path string) bool {
	return Config.Get(path).(bool)
}

//暂只支持顶层，不支持"A.B.C"这种path
func (c *ConfigType) Set(path string, value interface{}) {
	Config[path] = value
}

//暂只支持顶层，不支持"A.B.C"这种path
func Set(path string, value interface{}) {
	Config.Set(path, value)
}

func Init(customConf *ConfigType) {
	var defaultConf defaultConfig
	configFile := getConfigFile()
	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&defaultConf); err != nil {
			fmt.Println(err)
		}
	})

	if err := v.Unmarshal(&defaultConf); err != nil {
		fmt.Println(err)
	}

	funk.Map(funk.Keys(defaultConf), func(k string) string {
		if cv := customConf.Get(k); cv != nil {
			Set(k, cv)
		} else {
			Set(k, funk.Get(defaultConf, k))
		}
		return k
	})

	if GetBool("Debug") {
		j, _ := json.Marshal(Config)
		fmt.Printf("\nDebug模式：开启\n========Config========\n%v\n======================\n", string(j))
	}

}

func getConfigFile() string {
	// 配置文件的读取优先级: 命令行 > 环境变量 > 默认值
	var configFile string
	flag.StringVar(&configFile, "c", "", "choose config file.")
	flag.Parse()
	if configFile == "" {
		if configEnv := os.Getenv(ConfigEnv); configEnv == "release" {
			configFile = ProdConfigFile
			fmt.Printf("您正在使用生产环境")
		} else if configEnv == "test" {
			configFile = TestConfigFile
			fmt.Printf("您正在使用测试环境")
		} else {
			configFile = DevConfigFile
			fmt.Printf("您正在使用开发环境")
		}
	} else {
		fmt.Printf("您正在使用命令行的-c参数传递的值")
	}
	configFile = ConfigDir + configFile
	fmt.Printf(" config的路径为: %v\n", configFile)
	return configFile
}
