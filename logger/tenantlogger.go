package logger

import (
	"git.kuainiujinke.com/oa/oa-common-golang/config"
	"go.uber.org/zap"
)

var StoresLogger = make(map[string]*TenantLogger)

type TenantLogger struct {
	tenant   string
	loggerId string
	Zap      *zap.Logger
}

func NewTenantLogger() {
	logger := &TenantLogger{
		tenant:   config.PlatformAlias,
		loggerId: "",
		Zap:      NewZap(),
	}
	StoresLogger[config.PlatformAlias] = logger
}

func ByName(tenant, loggerId string) *TenantLogger {
	if _, ok := StoresLogger[tenant]; !ok {
		logger := StoresLogger[config.PlatformAlias]
		logger.tenant = tenant
		StoresLogger[tenant] = logger
	}

	StoresLogger[tenant].loggerId = loggerId
	return StoresLogger[tenant]
}

func (log *TenantLogger) Info(msg string, fields ...zap.Field) {
	msg = log.loggerId + " " + log.tenant + " " + msg
	log.Zap.Info(msg, fields...)
}

func (log *TenantLogger) Error(msg string, fields ...zap.Field) {
	msg = log.loggerId + " " + log.tenant + " " + msg
	log.Zap.Error(msg, fields...)
}
