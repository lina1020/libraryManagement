package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `json:"username" gorm:"uniqueIndex;not null"`
	PasswordHash string `json:"-" gorm:"column:password_hash"` // 不返回给前端
	Role         string `json:"role" gorm:"type:enum('user','admin');default:'user'"`
}
