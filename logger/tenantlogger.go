package logger

import (
	"go.uber.org/zap"
)

var StoresLogger = make(map[string]*TenantLogger)

type TenantLogger struct {
	tenant   string
	loggerId string
	Zap      *zap.Logger
}

func NewTenantLogger(loggerName string) {
	logger := &TenantLogger{
		tenant:   loggerName,
		loggerId: "",
		Zap:      NewZap(loggerName),
	}
	StoresLogger[loggerName] = logger
}

func ByName(tenant, loggerId string) *TenantLogger {
	if _, ok := StoresLogger[tenant]; !ok {
		NewTenantLogger(tenant)
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
