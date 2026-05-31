package controller

import (
	"BBingyan/internal/controller/param"
	"BBingyan/internal/model"
	"BBingyan/internal/utils"

	"github.com/labstack/echo/v4"
)

func Follow(c echo.Context) error {
	var req struct {
		Followee int `query:"followee"`
	}
	if err := c.Bind(&req); err != nil || req.Followee == 0 {
		return param.BadRequest(c, "followee required")
	}

	if err := model.CreateFollow(c.Request().Context(), utils.GetUID(c), req.Followee); err != nil {
		return param.InternalError(c, "")
	}
	return param.Success(c, nil)
}

func Unfollow(c echo.Context) error {
	var req struct {
		Followee int `query:"followee"`
	}
	if err := c.Bind(&req); err != nil || req.Followee == 0 {
		return param.BadRequest(c, "followee required")
	}

	if err := model.DeleteFollow(c.Request().Context(), utils.GetUID(c), req.Followee); err != nil {
		return param.InternalError(c, "")
	}
	return param.Success(c, nil)
}
