package utils

import (
	"os"
	"time"

	"BBingyan/internal/config"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger // 全局日志记录器

func InitLogger() {
	// 根据配置选择开发/生产模式
	var zapConfig zap.Config
	if config.Conf.Logger.Debug {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	// 日志分割配置
	// rotatelogs 会自动按时间切割日志文件
	writer, err := rotatelogs.New(
		config.Conf.Logger.Path+"%Y%m%d%H%M.log",
		rotatelogs.WithLinkName(config.Conf.Logger.Path+"latest.log"), // 软链接指向最新日志
		rotatelogs.WithMaxAge(7*24*time.Hour),                         // 日志保留 7 天
		rotatelogs.WithRotationTime(24*time.Hour),                     // 每 24 小时切割一次
	)
	if err != nil {
		panic("failed to init log rotatelogs: " + err.Error())
	}

	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 使用 ISO8601 时间格式

	level := zap.InfoLevel
	if config.Conf.Logger.Debug {
		level = zap.DebugLevel
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapConfig.EncoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(writer), zapcore.AddSync(os.Stdout)),
		level, // ← 生产环境用 Info，开发环境用 Debug
	)

	Logger = zap.New(core)
	Logger.Info("logger initialized successfully")
}
