package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar(16);comment:用户名;NOT NULL" json:"username"`
}
