package utils

import (
	"fmt"
	"time"

	"BBingyan/internal/config"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	UID        int `json:"uid"`
	Permission int `json:"perm"`
	jwt.RegisteredClaims
}

func GenerateToken(uid, perm int) (string, error) {
	claims := Claims{
		UID: uid, Permission: perm,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(
				time.Duration(config.Conf.Jwt.Expire) * time.Second,
			)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(config.Conf.Jwt.Secret))
}

func InitJWT(e *echo.Echo) {
	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(config.Conf.Jwt.Secret),
		TokenLookup: "header:Authorization:Bearer ",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &Claims{}
		},
		Skipper: func(c echo.Context) bool {
			path := c.Path()
			for _, p := range config.Conf.Jwt.SkippedPaths {
				if path == p {
					return true
				}
			}
			// Registration endpoint is open
			if c.Request().Method == "POST" && path == fmt.Sprintf("/%s/user", config.Conf.Server.Ver) {
				return true
			}
			return false
		},
	}))
}

// GetUID 从 echo-jwt 已验证的 token 中取 uid（不再二次 parse）
func GetUID(c echo.Context) int {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok || token == nil {
		return -1
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return -1
	}
	return claims.UID
}

// GetPermission 从 echo-jwt 已验证的 token 中取权限等级
func GetPermission(c echo.Context) int {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok || token == nil {
		return 0
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0
	}
	return claims.Permission
}
