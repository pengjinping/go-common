package logger

import (
	"go.uber.org/zap"
	"sync"
)

var StoresLogger = make(map[string]*TenantLogger)
var lock sync.Mutex

type TenantLogger struct {
	tenant   string
	loggerId string
	Zap      *zap.Logger
}

func NewTenantLogger(loggerName string) {
	if _, ok := StoresLogger[loggerName]; ok {
		return
	}

	lock.Lock()
	logger := &TenantLogger{
		tenant:   loggerName,
		loggerId: "",
		Zap:      NewZap(loggerName),
	}

	StoresLogger[loggerName] = logger
	lock.Unlock()
}

func ByName(tenant, loggerId string) *TenantLogger {
	NewTenantLogger(tenant)
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
