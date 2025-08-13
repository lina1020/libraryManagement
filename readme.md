# 简易图书管理系统技术文档

## 功能需求

提供 RESTful API 管理图书信息，使用go-gin框架完成服务端搭建，使用json序列化协议
* 能够将图书数据保存到 MySQL 数据库
* 能够检索图书列表、查看图书详细信息
* 包含完整的单元测试

> 不包含前端 UI 界面。

---

## 管理端接口说明

### 1. 添加书籍
- **方法**：`POST`
- **路径**：`/books/add`
- **描述**：新增一本图书记录
- **请求体示例**：
  ```json
  {
    "title": "Go语言编程",
    "count": 5,
    "isbn": "9787111111111"
  }
  ```
- **参数结构（BookInfoDTO）**：
  ```go
  type BookInfoDTO struct {
      Title string `json:"title" validate:"required"`
      Count uint   `json:"count" validate:"required"`
      ISBN  string `json:"isbn" validate:"required"`
  }
  ```

---

### 2. 删除书籍
- **方法**：`DELETE`
- **路径**：`/books/delete`
- **描述**：根据 ID 批量删除图书（软删除）
- **请求体示例**：
  ```json
  {
    "ids": ["1", "2", "3"]
  }
  ```
- **参数结构**：
  ```go
      IDs []string `json:"ids"`
  ```

---

### 3. 更新书籍
- **方法**：`PUT`
- **路径**：`/books/update`
- **描述**：更新指定图书信息，使用乐观锁防止并发冲突
- **请求体示例**：
  ```json
  {
    "id": 1,
    "title": "Go语言高级编程",
    "count": 10,
    "isbn": "9787111111111"
  }
  ```
- **参数结构（BookUpdateDTO）**：
  ```go
  type BookUpdateDTO struct {
      ID uint `json:"id" validate:"required"`
      BookInfoDTO
  }
  ```

---

### 4. 批量查询图书
- **方法**：`GET`
- **路径**：`/books/list`
- **描述**：根据可选条件模糊查询图书列表
- **支持查询参数**：
    - `title`：按书名模糊匹配
    - `isbn`：按 ISBN 精确匹配
- **示例请求**：
  ```
  GET /books/list?title=Go&isbn=9787111111111
  ```
- **参数结构（BookSearchDTO）**：
  ```go
  type BookSearchDTO struct {
      Title string `json:"title"`
      ISBN  string `json:"isbn"`
  }
  ```

---

## 数据库设计

```sql
CREATE TABLE books (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3),
    updated_at DATETIME(3),
    deleted_at DATETIME(3),
    title VARCHAR(16) NOT NULL COMMENT '书名',
    count BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '库存数量',
    isbn VARCHAR(13) NOT NULL COMMENT '国际标准书号',
    version INT NOT NULL DEFAULT 1 COMMENT '乐观锁版本号',

    UNIQUE KEY idx_isbn (isbn),
    INDEX idx_title (title),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 设计说明

| 字段 | 说明 |
|------|------|
| `deleted_at` | 实现软删除功能，保留删除记录 |
| `version` | 乐观锁控制字段，防止并发更新覆盖 |
| `idx_isbn` | 唯一索引，确保 ISBN 全局唯一 |
| `idx_title` | 普通索引，提升按书名检索性能 |

---

## 并发控制：乐观锁机制

在多管理员并发更新同一图书时，为避免数据丢失或覆盖问题，系统采用**数据库乐观锁**策略。

### 实现原理

1. 读取图书信息时携带当前 `version` 值
2. 更新请求中包含原始 `version`
3. 执行更新时使用条件更新语句

该机制在保证数据一致性的同时，避免了悲观锁带来的性能损耗。

---

## TODO 待办事项

- [ ] **用户与管理员认证模块**
    - 用户/管理员注册
    - 登录鉴权（JWT 或 Session）
- [ ] **用户端功能开发**
    - 图书借阅功能
    - 图书归还功能
    - 用户可查询图书列表与详情
