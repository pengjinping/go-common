package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	ConfigEnv      = "GIN_MODE"
	ConfigDir      = "./deployment/"
	DevConfigFile  = "application-dev.yaml"
	TestConfigFile = "application-test.yaml"
	ProdConfigFile = "application-prod.yaml"
)

var config = viper.GetViper()

// IsSet checks to see if the key has been set in any of the data locations.
// IsSet is case-insensitive for a key.
func IsSet(key string) bool { return config.IsSet(key) }

// Set sets the value for the key in the override register.
// Set is case-insensitive for a key.
// Will be used instead of values obtained via
// flags, config file, ENV, default, or key/value store.
func Set(key string, value interface{}) { config.Set(key, value) }

// Get can retrieve any value given the key to use.
// Get is case-insensitive for a key.
// Get has the behavior of returning the value associated with the first
// place from where it is set. Viper will check in the following order:
// override, flag, env, config file, key/value store, default
//
// Get returns an interface. For a specific value use one of the Get____ methods.
func Get(key string) interface{} { return config.Get(key) }

// GetString returns the value associated with the key as a string.
func GetString(key string) string { return config.GetString(key) }

// GetBool returns the value associated with the key as a boolean.
func GetBool(key string) bool { return config.GetBool(key) }

// GetInt returns the value associated with the key as an integer.
func GetInt(key string) int { return config.GetInt(key) }

// GetInt32 returns the value associated with the key as an integer.
func GetInt32(key string) int32 { return config.GetInt32(key) }

// GetInt64 returns the value associated with the key as an integer.
func GetInt64(key string) int64 { return config.GetInt64(key) }

// GetUint returns the value associated with the key as an unsigned integer.
func GetUint(key string) uint { return config.GetUint(key) }

// GetUint32 returns the value associated with the key as an unsigned integer.
func GetUint32(key string) uint32 { return config.GetUint32(key) }

// GetUint64 returns the value associated with the key as an unsigned integer.
func GetUint64(key string) uint64 { return config.GetUint64(key) }

// GetFloat64 returns the value associated with the key as a float64.
func GetFloat64(key string) float64 { return config.GetFloat64(key) }

// GetTime returns the value associated with the key as time.
func GetTime(key string) time.Time { return config.GetTime(key) }

// GetDuration returns the value associated with the key as a duration.
func GetDuration(key string) time.Duration { return config.GetDuration(key) }

// GetIntSlice returns the value associated with the key as a slice of int values.
func GetIntSlice(key string) []int { return config.GetIntSlice(key) }

// GetStringSlice returns the value associated with the key as a slice of strings.
func GetStringSlice(key string) []string { return config.GetStringSlice(key) }

// GetStringMap returns the value associated with the key as a map of interfaces.
func GetStringMap(key string) map[string]interface{} { return config.GetStringMap(key) }

// GetStringMapString returns the value associated with the key as a map of strings.
func GetStringMapString(key string) map[string]string { return config.GetStringMapString(key) }

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func GetStringMapStringSlice(key string) map[string][]string {
	return config.GetStringMapStringSlice(key)
}

// GetSizeInBytes returns the size of the value associated with the given key
// in bytes.
func GetSizeInBytes(key string) uint { return config.GetSizeInBytes(key) }

// UnmarshalKey takes a single key and unmarshals it into a Struct.
func UnmarshalKey(key string, rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	return config.UnmarshalKey(key, rawVal, opts...)
}

// GetViper gets the global Viper instance.
func GetViper() *viper.Viper {
	return config
}

func Init(customConf *ConfigType) {
	configFile := getConfigFile()
	config.SetConfigFile(configFile)
	config.SetConfigType("yaml")
	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	config.WatchConfig()

	// config.SetEnvPrefix("OA")
	config.AutomaticEnv()

	config.MergeConfigMap(*customConf)

	if GetBool("Debug") {
		j, _ := json.Marshal(config.AllSettings())
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
