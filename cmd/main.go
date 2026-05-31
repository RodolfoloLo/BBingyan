package main

import (
	"BBingyan/internal/config"
	"BBingyan/internal/model"
	"BBingyan/internal/router"
	"BBingyan/internal/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	config.InitConfig()
	e := echo.New()

	utils.InitLogger(e)
	// 确保日志缓冲区的日志被写出
	defer utils.Logger.Sync()

	model.InitDB()
	model.InitDefaultAdmin()
	utils.InitRedis()
	utils.InitJWT(e)
	router.InitRouter(e)

	if err := e.Start(":" + config.Conf.Server.Port); err != nil {
		utils.Logger.Fatal("server start failed", zap.Error(err))
	}
}
