package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TenantLogger struct {
	Zap *zap.Logger
}

func (tl *TenantLogger) Info(context *gin.Context, msg string, fields ...zap.Field) {
	host := context.Request.Host
	tl.Zap.Info(host+msg, fields...)
}
func (tl *TenantLogger) Error(context *gin.Context, msg string, fields ...zap.Field) {
	host := context.Request.Host
	tl.Zap.Error(host+msg, fields...)
}
