package controller

import (
	"BBingyan/internal/controller/param"
	"BBingyan/internal/model"
	"BBingyan/internal/utils"

	"github.com/labstack/echo/v4"
)

func Like(c echo.Context) error {
	var req param.PIDQuery
	if err := c.Bind(&req); err != nil || req.PID == 0 {
		return param.BadRequest(c, "pid required")
	}

	if err := model.CreateLike(c.Request().Context(), utils.GetUID(c), req.PID); err != nil {
		return param.InternalError(c, "")
	}
	return param.Success(c, nil)
}

func Unlike(c echo.Context) error {
	var req param.PIDQuery
	if err := c.Bind(&req); err != nil || req.PID == 0 {
		return param.BadRequest(c, "pid required")
	}

	if err := model.DeleteLike(c.Request().Context(), utils.GetUID(c), req.PID); err != nil {
		return param.InternalError(c, "")
	}
	return param.Success(c, nil)
}
