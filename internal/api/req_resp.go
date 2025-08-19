package api

type BookInfoReq struct {
	Title string `json:"title" validate:"required"`
	Count uint   `json:"count" validate:"required"`
	ISBN  string `json:"isbn" validate:"required"`

	Author  string `json:"author"`
	Content string `json:"content"`
	Summary string `json:"summary"`
}

type BookUpdateReq struct {
	ID uint `json:"id" validate:"required"`
	BookInfoReq
}

type BookSearchReq struct {
	Title string `json:"title"`
	ISBN  string `json:"isbn"`

	Author   string `json:"author"`
	Content  string `json:"content"`   // 模糊搜索内容
	Keyword  string `json:"keyword"`   // 全文搜索关键词
	Page     int    `json:"page"`      // 分页页码
	PageSize int    `json:"page_size"` // 每页大小
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

	Author  string `json:"author"`
	Content string `json:"content,omitempty"` // 列表查询时可能不返回内容
	Summary string `json:"summary"`
}

type BookSearchResp struct {
	Books      []BookInfoResp `json:"books"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

type LoginResp struct {
	Token  string `json:"token"`
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}
