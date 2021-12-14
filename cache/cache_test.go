package cache

import "testing"

func TestRegister(t *testing.T) {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Register(tt.args.driver)
		})
	}
}
