package param

import "github.com/labstack/echo/v4"

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

type PageResp struct {
	Items    any   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

type Paging struct {
	ID       int `query:"id"`
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
	Sort     int `query:"sort"` // 0=time, 1=comments, 2=likes
}

func (p *Paging) Default() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 || p.PageSize > 100 {
		p.PageSize = 20
	}
}

func (p Paging) SortClause() string {
	switch p.Sort {
	case 1:
		return "comments desc"
	case 2:
		return "likes desc"
	default:
		return "created_at desc"
	}
}

func Success(c echo.Context, data any) error {
	return c.JSON(200, Resp{Code: 200, Msg: "ok", Data: data})
}

func BadRequest(c echo.Context, msg string) error {
	if msg == "" {
		msg = "bad request"
	}
	return c.JSON(400, Resp{Code: 400, Msg: msg})
}

func Unauthorized(c echo.Context, msg string) error {
	if msg == "" {
		msg = "unauthorized"
	}
	return c.JSON(401, Resp{Code: 401, Msg: msg})
}

func Forbidden(c echo.Context, msg string) error {
	if msg == "" {
		msg = "forbidden"
	}
	return c.JSON(403, Resp{Code: 403, Msg: msg})
}

func NotFound(c echo.Context, msg string) error {
	if msg == "" {
		msg = "not found"
	}
	return c.JSON(404, Resp{Code: 404, Msg: msg})
}

func Conflict(c echo.Context, msg string) error {
	if msg == "" {
		msg = "conflict"
	}
	return c.JSON(409, Resp{Code: 409, Msg: msg})
}

func TooMany(c echo.Context, retryAfter string) error {
	return c.JSON(429, Resp{Code: 429, Msg: "too frequent, retry after: " + retryAfter})
}

func InternalError(c echo.Context, msg string) error {
	if msg == "" {
		msg = "internal error"
	}
	return c.JSON(500, Resp{Code: 500, Msg: msg})
}
