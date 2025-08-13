db:library

table：
book:
Title string `gorm:"column:title;type:varchar(16);comment:书名;NOT NULL" json:"title" json:"title,omitempty"`
Count uint   `gorm:"column:count;type:bigint unsigned;default:0;comment:数量;NOT NULL" json:"count"`
ISBN  string `gorm:"column:isbn;type:varchar(13);uniqueIndex:idx_isbn;comment:编码;NOT NULL" json:"isbn"`



user:
Username string `gorm:"column:username;type:varchar(16);comment:用户名;NOT NULL" json:"username"`
password
userid