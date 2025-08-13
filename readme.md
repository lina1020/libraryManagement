# 简易图书管理系统

### 功能需求

提供 RESTful API 管理图书信息，使用go-gin框架完成服务端搭建，使用json序列化协议
能够将图书数据保存到 MySQL 数据库
能够检索图书列表、查看图书详细信息
包含完整的单元测试

_不需要实现UI界面_


### 管理端接口说明

1.添加书籍 POST： /books/add
    接受参数：
        BookInfoDTO:
            Title string `json:"title" validate:"required"`
            Count uint  `json:"count" validate:"required"`
            ISBN  string `json:"isbn" validate:"required"`

2.删除书籍 DELETE: /books/delete
    接受参数：
        ids []string

3.更新书籍 PUT: /books/update
    接受参数:
        BookUpdateDTO:
            ID uint `json:"id" validate:"required"`
            BookInfoDTO

4.批量查询 GET： /books/list
    接受参数:
        BookSearchDTO:
            Title string `json:"title"`
            ISBN  string `json:"isbn"`

### 数据库设计

id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
created_at DATETIME(3),
updated_at DATETIME(3),
deleted_at DATETIME(3),
title VARCHAR(16) NOT NULL COMMENT '书名',
count BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '数量',
isbn VARCHAR(13) NOT NULL COMMENT '编码',
version INT NOT NULL DEFAULT 1 COMMENT '乐观锁版本号',

UNIQUE KEY idx_isbn (isbn),
INDEX idx_title (title),
INDEX idx_deleted_at (deleted_at)

version 字段用于实现乐观锁机制，防止并发更新导致的数据覆盖问题。
建立唯一索引 idx_isbn 保证 ISBN 的全局唯一性。
普通索引 idx_title 提升按书名检索的性能。

### 并发问题
在多管理员并发更新同一本书籍的场景下，为避免数据丢失或覆盖问题，采用数据库乐观锁策略进行控制。

### TODO
1. [ ] 用户、管理员注册登入实现
2. [ ] 用户端借阅、归还、查询功能实现