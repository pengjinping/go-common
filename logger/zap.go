package logger

import (
	"fmt"
	"os"
	"time"

	"git.kuainiujinke.com/oa/oa-common-golang/config"
	"git.kuainiujinke.com/oa/oa-common-golang/utils"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZap(tenant string) *zap.Logger {
	zapConf := getZapConfig()
	dirName := zapConf.Directory + "/" + tenant
	if ok, _ := utils.PathExists(dirName); !ok { // 判断是否有Directory文件夹
		fmt.Printf("create %v directory\n", dirName)
		_ = os.Mkdir(dirName, os.ModePerm)
	}
	// 调试级别
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.DebugLevel
	})
	// 日志级别
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.InfoLevel
	})
	// 警告级别
	warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.WarnLevel
	})
	// 错误级别
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})

	cores := [...]zapcore.Core{
		getEncoderCore(fmt.Sprintf("./%s/debug.log", dirName), debugPriority),
		getEncoderCore(fmt.Sprintf("./%s/info.log", dirName), infoPriority),
		getEncoderCore(fmt.Sprintf("./%s/warn.log", dirName), warnPriority),
		getEncoderCore(fmt.Sprintf("./%s/error.log", dirName), errorPriority),
	}
	logger := zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller(), zap.AddCallerSkip(2))

	if zapConf.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}

	return logger
}

func getZapConfig() config.ZapConfig {
	var zapConf config.ZapConfig
	if err := config.UnmarshalKey("Zap", &zapConf); err != nil {
		fmt.Printf("Zap config init failed: %s\n", err)
	}
	return zapConf
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig() (config zapcore.EncoderConfig) {
	zapConf := getZapConfig()
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  zapConf.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case zapConf.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case zapConf.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case zapConf.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case zapConf.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	if getZapConfig().Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore(fileName string, level zapcore.LevelEnabler) (core zapcore.Core) {
	writer := GetWriteSyncEr(fileName) // 使用file-rotatelogs进行日志分割
	return zapcore.NewCore(getEncoder(), writer, level)
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(getZapConfig().Prefix + "2006/01/02 - 15:04:05.000"))
}
func GetWriteSyncEr(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file, //日志文件的位置
		MaxSize:    10,   //在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 200,  //保留旧文件的最大个数
		MaxAge:     30,   //保留旧文件的最大天数
		Compress:   true, //是否压缩/归档旧文件
	}

	if getZapConfig().LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	return zapcore.AddSync(lumberJackLogger)
}
