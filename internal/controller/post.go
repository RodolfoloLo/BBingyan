package controller

import (
	"BBingyan/internal/controller/param"
	"BBingyan/internal/model"
	"BBingyan/internal/utils"

	"github.com/labstack/echo/v4"
)

func CreatePost(c echo.Context) error {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		NID     int    `json:"nid"`
	}
	if err := c.Bind(&req); err != nil {
		return param.BadRequest(c, "")
	}

	post := &model.Post{
		UID:   utils.GetUID(c),
		Title: req.Title,
		NID:   req.NID,
	}
	if req.Content != "" {
		post.Content = &req.Content
	}

	result, err := model.CreatePost(c.Request().Context(), post)
	if err != nil {
		return param.InternalError(c, "")
	}
	return param.Success(c, result)
}

func ListPosts(c echo.Context) error {
	var p param.Paging
	c.Bind(&p)
	p.Default()

	posts, total, err := model.GetPosts(c.Request().Context(), p)
	if err != nil {
		return param.InternalError(c, "")
	}
	return param.Success(c, param.PageResp{
		Items:    posts,
		Total:    total,
		Page:     p.Page,
		PageSize: p.PageSize,
	})
}

func GetPost(c echo.Context) error {
	var req struct {
		PID int `query:"pid"`
	}
	if err := c.Bind(&req); err != nil || req.PID == 0 {
		return param.BadRequest(c, "pid required")
	}

	post, err := model.GetPostByPID(c.Request().Context(), req.PID)
	if err != nil {
		return param.NotFound(c, "")
	}
	return param.Success(c, post)
}

func DeletePost(c echo.Context) error {
	var req struct {
		PID int `query:"pid"`
	}
	if err := c.Bind(&req); err != nil || req.PID == 0 {
		return param.BadRequest(c, "pid required")
	}

	post, err := model.GetPostByPID(c.Request().Context(), req.PID)
	if err != nil {
		return param.NotFound(c, "")
	}

	uid := utils.GetUID(c)
	if utils.GetPermission(c) < 1 && post.UID != uid {
		return param.Forbidden(c, "not your post")
	}

	if err := model.DeletePost(c.Request().Context(), req.PID); err != nil {
		return param.InternalError(c, "")
	}
	return param.Success(c, nil)
}
