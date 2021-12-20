package config

// 当前服务是否在生产环境
func IsProduction() bool {
	// config.GetString("env"): prod, test, dev
	return config.GetString("env") == "prod"
}
