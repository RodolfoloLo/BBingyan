package router

import (
	"BBingyan/internal/config"
	"BBingyan/internal/controller"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRouter(e *echo.Echo) {
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	v := config.Conf.Server.Ver

	// 白名单 — 不需要登录
	e.POST(fmt.Sprintf("/%s/user/token", v), controller.Login)
	e.POST(fmt.Sprintf("/%s/user", v), controller.Register)
	e.GET(fmt.Sprintf("/%s/verify", v), controller.SendValidation)

	// 需要 JWT
	e.GET(fmt.Sprintf("/%s/user", v), controller.GetUser)
	e.DELETE(fmt.Sprintf("/%s/user", v), controller.DeleteUser)

	e.POST(fmt.Sprintf("/%s/post", v), controller.CreatePost)
	e.GET(fmt.Sprintf("/%s/post", v), controller.ListPosts)
	e.GET(fmt.Sprintf("/%s/post/pid", v), controller.GetPost)
	e.DELETE(fmt.Sprintf("/%s/post", v), controller.DeletePost)

	e.POST(fmt.Sprintf("/%s/comment", v), controller.CreateComment)
	e.GET(fmt.Sprintf("/%s/comment/pid", v), controller.ListCommentsByPID)
	e.DELETE(fmt.Sprintf("/%s/comment", v), controller.DeleteComment)

	e.POST(fmt.Sprintf("/%s/like", v), controller.Like)
	e.DELETE(fmt.Sprintf("/%s/like", v), controller.Unlike)

	e.POST(fmt.Sprintf("/%s/follow", v), controller.Follow)
	e.DELETE(fmt.Sprintf("/%s/follow", v), controller.Unfollow)

	e.POST(fmt.Sprintf("/%s/node", v), controller.CreateNode)
	e.GET(fmt.Sprintf("/%s/node", v), controller.ListNodes)
	e.DELETE(fmt.Sprintf("/%s/node", v), controller.DeleteNode)
}
