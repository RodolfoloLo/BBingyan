package utils

import (
	"os"
	"time"

	"BBingyan/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger(e *echo.Echo) {
	// 根据 debug 配置决定编码器：开发用 console（可读），生产用 JSON
	var encoder zapcore.Encoder
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if config.Conf.Logger.Debug {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		prodConfig := zap.NewProductionEncoderConfig()
		prodConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewJSONEncoder(prodConfig)
	}

	// 文件输出 — rotatelogs 按天切割，保留 7 天
	fileWriter, err := rotatelogs.New(
		config.Conf.Logger.Path+"%Y%m%d.log",
		rotatelogs.WithLinkName(config.Conf.Logger.Path+"latest.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	var core zapcore.Core
	if err != nil {
		// 文件日志不可用就降级，只输出到终端
		core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
		Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		Logger.Warn("rotatelogs failed, using stdout only", zap.Error(err))
	} else {
		level := zapcore.InfoLevel
		if config.Conf.Logger.Debug {
			level = zapcore.DebugLevel
		}
		core = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(fileWriter),
			zapcore.AddSync(os.Stdout),
		), level)
		Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	}

	Logger.Info("logger ready")

	// HTTP 请求日志中间件
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogMethod:   true,
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			Logger.Info("request",
				zap.String("method", v.Method),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
				zap.String("ip", v.RemoteIP),
			)
			return nil
		},
	}))
}
