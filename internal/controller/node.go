package controller

import (
	"BBingyan/internal/controller/param"
	"BBingyan/internal/model"
	"BBingyan/internal/utils"

	"github.com/labstack/echo/v4"
)

func CreateNode(c echo.Context) error {
	if utils.GetPermission(c) < 1 {
		return param.Forbidden(c, "admin only")
	}

	var req param.CreateNodeRequest
	if err := c.Bind(&req); err != nil || req.Name == "" {
		return param.BadRequest(c, "")
	}

	node := &model.Node{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := model.CreateNode(c.Request().Context(), node); err != nil {
		return param.InternalError(c, "")
	}
	return param.Success(c, node)
}

func ListNodes(c echo.Context) error {
	nodes, err := model.ListNodes(c.Request().Context())
	if err != nil {
		return param.InternalError(c, "")
	}
	return param.Success(c, nodes)
}

func DeleteNode(c echo.Context) error {
	if utils.GetPermission(c) < 1 {
		return param.Forbidden(c, "admin only")
	}

	var req param.NIDQuery
	if err := c.Bind(&req); err != nil || req.NID == 0 {
		return param.BadRequest(c, "nid required")
	}

	if err := model.DeleteNode(c.Request().Context(), req.NID); err != nil {
		return param.InternalError(c, "")
	}
	return param.Success(c, nil)
}
