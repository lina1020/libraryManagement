package api

type BookInfoReq struct {
	Title string `json:"title" validate:"required"`
	Count uint   `json:"count" validate:"required"`
	ISBN  string `json:"isbn" validate:"required"`
}

type BookUpdateReq struct {
	ID uint `json:"id" validate:"required"`
	BookInfoReq
}

type BookSearchReq struct {
	Title string `json:"title"`
	ISBN  string `json:"isbn"`
}

type RegisterReq struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"omitempty,oneof=user admin"` // 可选，默认为 user
}

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type BookInfoResp struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Count uint   `json:"count"`
	ISBN  string `json:"isbn"`
}

type LoginResp struct {
	Token  string `json:"token"`
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}
