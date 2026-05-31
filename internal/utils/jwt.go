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

func ParseToken(raw string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(raw, claims, func(t *jwt.Token) (any, error) {
		return []byte(config.Conf.Jwt.Secret), nil
	})
	return claims, err
}

func InitJWT(e *echo.Echo) {
	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(config.Conf.Jwt.Secret),
		TokenLookup: "header:Authorization:Bearer ",
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

func GetUID(c echo.Context) int {
	token, _ := c.Get("user").(*jwt.Token)
	if token == nil {
		return -1
	}
	claims, err := ParseToken(token.Raw)
	if err != nil {
		return -1
	}
	return claims.UID
}

func GetPermission(c echo.Context) int {
	token, _ := c.Get("user").(*jwt.Token)
	if token == nil {
		return 0
	}
	claims, err := ParseToken(token.Raw)
	if err != nil {
		return 0
	}
	return claims.Permission
}
