package controller

import (
	"BBingyan/internal/controller/param"
	"BBingyan/internal/service"
	"BBingyan/internal/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func SendValidation(c echo.Context) error {
	email := c.QueryParam("mail")
	if email == "" {
		return param.BadRequest(c, "mail required")
	}

	ok, wait, err := utils.CanResend(c.Request().Context(), email)
	if err != nil {
		return param.InternalError(c, "")
	}
	if !ok {
		return param.TooMany(c, wait.String())
	}

	code := utils.GenerateCode()
	if err := utils.SetCode(c.Request().Context(), email, code); err != nil {
		return param.InternalError(c, "")
	}

	// 异步发邮件，不阻塞 HTTP 响应
	go func() {
		if err := service.SendValidationCode(email, code); err != nil {
			utils.Logger.Error("send email failed", zap.String("to", email), zap.Error(err))
		}
	}()

	return param.Success(c, nil)
}
