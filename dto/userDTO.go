package dto

// RegisterDTO 注册请求
type RegisterDTO struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"omitempty,oneof=user admin"` // 可选，默认为 user
}

// LoginDTO 登录请求
type LoginDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
