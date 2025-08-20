package model

import "gorm.io/gorm"

//GORM 默认会将结构体名转换为复数形式并将其用作表名。如果你定义的 Go 结构体名称是 Book，
//那么 GORM 会自动匹配到名为 books 的数据库表，因为它是根据结构体名 Book 转换成复数形式来决定表名的。

type Book struct {
	gorm.Model
	Title   string `gorm:"column:title;type:varchar(255);comment:书名;NOT NULL" json:"title" json:"title,omitempty"`
	Count   uint   `gorm:"column:count;type:bigint unsigned;default:0;comment:数量;NOT NULL" json:"count"`
	ISBN    string `gorm:"column:isbn;type:varchar(17);uniqueIndex:idx_isbn;comment:编码;NOT NULL" json:"isbn"`
	Author  string `gorm:"column:author;type:varchar(100);comment:作者" json:"author"`
	Content string `gorm:"column:content;type:longtext;comment:书本内容" json:"content"`
	Summary string `gorm:"column:summary;type:text;comment:内容摘要" json:"summary"`

	Version int `gorm:"column:version;default:1;comment:乐观锁版本号" json:"version"`
}

// ESBookDocument ES中的书籍文档结构
type ESBookDocument struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Count   uint   `json:"count"`
	Author  string `json:"author"`
	ISBN    string `json:"isbn"`
	Content string `json:"content"`
	Summary string `json:"summary"`
}
