package controller

import (
	"errors"

	"BBingyan/internal/config"
	"BBingyan/internal/controller/param"
	"BBingyan/internal/model"
	"BBingyan/internal/utils"

	"net/mail"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func Register(c echo.Context) error {
	var req param.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return param.BadRequest(c, "")
	}

	// 验证邮箱格式
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return param.BadRequest(c, "invalid email")
	}

	code := c.QueryParam("code")
	valid, err := utils.ValidateCode(c.Request().Context(), req.Email, code)
	if err != nil || !valid {
		return param.Unauthorized(c, "invalid code")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Logger.Error("bcrypt failed", zap.Error(err))
		return param.InternalError(c, "")
	}
	user := &model.User{
		Username: req.Username,
		Password: string(hash),
		Email:    req.Email,
	}
	if err := model.CreateUser(c.Request().Context(), user); err != nil {
		if errors.Is(err, model.ErrUserAlreadyExist) {
			return param.Conflict(c, "username taken")
		}
		return param.InternalError(c, "")
	}
	return param.Success(c, nil)
}

func Login(c echo.Context) error {
	var req param.LoginRequest
	if err := c.Bind(&req); err != nil {
		return param.BadRequest(c, "")
	}

	user, err := model.GetUserByUsername(c.Request().Context(), req.Username)
	if errors.Is(err, model.ErrUserNotFound) {
		return param.Unauthorized(c, "wrong username or password")
	}
	if err != nil {
		return param.InternalError(c, "")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return param.Unauthorized(c, "wrong username or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Permission)
	if err != nil {
		utils.Logger.Error("generate token failed", zap.Error(err))
		return param.InternalError(c, "")
	}
	return param.Success(c, map[string]any{
		"token":      token,
		"expires_in": config.Conf.Jwt.Expire,
	})
}

func GetUser(c echo.Context) error {
	var req param.UserQuery
	c.Bind(&req)

	var user *model.User
	var err error
	switch {
	case req.ID != 0:
		user, err = model.GetUserByID(c.Request().Context(), req.ID)
	case req.Username != "":
		user, err = model.GetUserByUsername(c.Request().Context(), req.Username)
	default:
		user, err = model.GetUserByID(c.Request().Context(), utils.GetUID(c))
	}
	if errors.Is(err, model.ErrUserNotFound) {
		return param.NotFound(c, "")
	}
	if err != nil {
		return param.InternalError(c, "")
	}
	return param.Success(c, user)
}

func DeleteUser(c echo.Context) error {
	if utils.GetPermission(c) < 1 {
		return param.Forbidden(c, "admin only")
	}
	var req param.IDQuery
	c.Bind(&req)
	if err := model.DeleteUser(c.Request().Context(), req.ID); err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return param.NotFound(c, "")
		}
		return param.InternalError(c, "")
	}
	return param.Success(c, nil)
}
