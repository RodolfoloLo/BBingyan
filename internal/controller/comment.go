package controller

import (
	"BBingyan/internal/controller/param"
	"BBingyan/internal/model"
	"BBingyan/internal/utils"

	"github.com/labstack/echo/v4"
)

func CreateComment(c echo.Context) error {
	var req struct {
		PID     int    `json:"pid"`
		Content string `json:"content"`
	}
	if err := c.Bind(&req); err != nil || req.PID == 0 || req.Content == "" {
		return param.BadRequest(c, "")
	}

	comment := &model.Comment{
		UID:     utils.GetUID(c),
		PID:     req.PID,
		Content: req.Content,
	}
	if err := model.CreateComment(c.Request().Context(), comment); err != nil {
		return param.InternalError(c, "")
	}
	return param.Success(c, comment)
}

func ListCommentsByPID(c echo.Context) error {
	var req struct {
		PID int `query:"pid"`
	}
	if err := c.Bind(&req); err != nil || req.PID == 0 {
		return param.BadRequest(c, "pid required")
	}

	comments, err := model.ListCommentsByPID(c.Request().Context(), req.PID)
	if err != nil {
		return param.InternalError(c, "")
	}
	return param.Success(c, comments)
}

func DeleteComment(c echo.Context) error {
	var req struct {
		CID int `query:"cid"`
	}
	if err := c.Bind(&req); err != nil || req.CID == 0 {
		return param.BadRequest(c, "cid required")
	}

	comment, err := model.GetCommentByCID(c.Request().Context(), req.CID)
	if err != nil {
		return param.NotFound(c, "")
	}

	uid := utils.GetUID(c)
	if utils.GetPermission(c) < 1 && comment.UID != uid {
		return param.Forbidden(c, "not your comment")
	}

	if err := model.DeleteComment(c.Request().Context(), req.CID); err != nil {
		return param.InternalError(c, "")
	}
	return param.Success(c, nil)
}
