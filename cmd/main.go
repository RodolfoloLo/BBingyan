package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"BBingyan/config"
	"BBingyan/internal/model"
	"BBingyan/internal/router"
	"BBingyan/internal/utils"
)

func main() {
	config.InitConfig()
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	utils.InitLogger()
	defer utils.Logger.Sync()

	model.InitDB()
	model.InitDefaultAdmin()
	utils.InitRedis()
	utils.InitJWT(e)
	router.InitRouter(e)

	go func() {
		addr := ":" + config.Conf.Server.Port
		utils.Logger.Info("listening", zap.String("addr", addr))
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down: " + err.Error())
		}
	}()

	// 等 SIGINT/SIGTERM, 最多 10s 清理
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Logger.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
