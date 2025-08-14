package vo

// LoginVO LoginResponse 登录响应
type LoginVO struct {
	Token  string `json:"token"`
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}
