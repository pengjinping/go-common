package logger

import (
	"context"
	"git.kuainiujinke.com/oa/oa-common-golang/config"
	"go.uber.org/zap"
)

func Init() {
	NewTenantLogger()
}

func Get(ctx context.Context) *TenantLogger {
	loggerId := ""
	if loggerIdName := ctx.Value("loggerId"); loggerIdName != nil {
		loggerId = loggerIdName.(string)
	}

	if loggerName := ctx.Value("logger"); loggerName != nil {
		logger := loggerName.(*TenantLogger)
		logger.loggerId = loggerId
		return logger
	}

	tenant := config.PlatformAlias
	if tenantName := ctx.Value("tenant"); tenantName != nil {
		tenant = tenantName.(string)
	}

	return ByName(tenant, loggerId)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	Get(ctx).Info(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	Get(ctx).Error(msg, fields...)
}
