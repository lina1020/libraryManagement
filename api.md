## 一、管理端接口说明

### 1. 添加书籍
- **方法**：`POST`
- **路径**：`/admin/books/add`
- **权限**：管理员（`admin`）
- **描述**：新增一本图书
- **请求体**：
  ```json
  {
    "title": "Go语言编程",
    "count": 5,
    "isbn": "9787111111111"
  }
  ```
- **结构体**：
  ```go
  type BookInfoReq struct {
      Title string `json:"title" validate:"required"`
      Count uint   `json:"count" validate:"required"`
      ISBN  string `json:"isbn" validate:"required"`
  }
  ```

---

### 2. 删除书籍
- **方法**：`DELETE`
- **路径**：`/admin/books/delete`
- **权限**：管理员
- **描述**：根据 ID 批量软删除图书
- **请求体**：
  ```json
  { "ids": ["1", "2"] }
  ```

---

### 3. 更新书籍
- **方法**：`PUT`
- **路径**：`/admin/books/update`
- **权限**：管理员
- **描述**：更新图书信息（带乐观锁）
- **请求体**：
  ```json
  {
    "id": 1,
    "title": "Go语言高级编程",
    "count": 8,
    "isbn": "9787111111111"
  }
  ```

---

### 4. 批量查询图书
- **方法**：`GET`
- **路径**：`/api/books/list`
- **权限**：所有登录用户
- **描述**：按书名或 ISBN 查询图书列表
- **参数示例**：
  ```
  GET /api/books/list?title=Go&isbn=9787111111111
  ```

---

## 二、用户认证接口

### 1. 用户注册
- **方法**：`POST`
- **路径**：`/auth/register`
- **权限**：公开
- **描述**：注册新用户或管理员
- **请求体**：
  ```json
  {
    "username": "alice",
    "password": "123456",
    "role": "user" // 可选，默认 user
  }
  ```

---

### 2. 用户登录
- **方法**：`POST`
- **路径**：`/auth/login`
- **权限**：公开
- **描述**：登录并获取 JWT Token
- **请求体**：
  ```json
  {
    "username": "alice",
    "password": "123456"
  }
  ```
- **响应**：
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.xxxxx",
    "user_id": 1,
    "role": "user"
  }
  ```
