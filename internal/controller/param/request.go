package param

// ——— Query param structs ———

type PIDQuery struct {
	PID int `query:"pid"`
}

type CIDQuery struct {
	CID int `query:"cid"`
}

type NIDQuery struct {
	NID int `query:"nid"`
}

type IDQuery struct {
	ID int `query:"id"`
}

type FolloweeQuery struct {
	Followee int `query:"followee"`
}

type UserQuery struct {
	ID       int    `query:"id"`
	Username string `query:"username"`
}

// ——— JSON body structs ———

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	NID     int    `json:"nid"`
}

type CreateCommentRequest struct {
	PID     int    `json:"pid"`
	Content string `json:"content"`
}

type CreateNodeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
