package cache

import (
	"git.kuainiujinke.com/oa/oa-common-golang/config"
	"testing"
)

var redisConfig = config.RedisConfig{
	Host:     "127.0.0.1",
	Port:     6379,
	Password: "123456",
	DBName:   0,
}

func Test_Register(t *testing.T) {
	type args struct {
		driver string
	}
	tests := []struct {
		name    string
		args    args
	}{
		{
			name: "memory",
			args: args{
				driver: "memory",
			},
		},
		{
			name: "redis",
			args: args{
				driver: "redis",
			},
		},
	}

	config.Set("Redis", redisConfig)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			register(tt.args.driver)
		})
	}
}

func Test_configDriver(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "测试默认驱动",
			want: "memory",
		},
		{
			name: "设置驱动Redis",
			want: "redis",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want == "redis" {
				config.Set("Cache", config.CacheConfig{
					Driver: "redis",
				})
			}

			if got := configDriver(); got != tt.want {
				t.Errorf("configDriver() = %v, want %v", got, tt.want)
			}
		})
	}
}