# 简易图书管理系统

## 一、功能需求

提供 RESTful API 管理图书信息，使用go-gin框架完成服务端搭建，使用json序列化协议
* 能够将图书数据保存到 MySQL 数据库
* 能够检索图书列表、查看图书详细信息
* 包含完整的单元测试

> 不包含前端 UI 界面。

---
## 二、功能模块

| 模块 | 功能 |
|------|------|
| 图书管理 | 增删改查、批量查询、乐观锁控制 |
| 用户系统 | 用户/管理员注册、登录、角色权限控制 |
| 安全认证 | JWT 鉴权、密码加密存储、接口权限校验 |
| 可测试性 | 核心逻辑覆盖单元测试 |

单元测试覆盖率：
* ok      LibraryManagement/controller    0.481s  coverage: 66.7% of statements
* ok      LibraryManagement/dao   0.697s  coverage: 61.0% of statements
* ok      LibraryManagement/service       1.146s  coverage: 38.5% of statements
* ok      LibraryManagement/utils 0.807s  coverage: 90.0% of statements

---

## 三、管理端接口说明

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

## 四、用户认证接口

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

---

## 五、数据库设计

### 1. 图书表 `books`

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

#### 设计说明

| 字段 | 说明 |
|------|------|
| `deleted_at` | 实现软删除功能，保留删除记录 |
| `version` | 乐观锁控制字段，防止并发更新覆盖 |
| `idx_isbn` | 唯一索引，确保 ISBN 全局唯一 |
| `idx_title` | 普通索引，提升按书名检索性能 |

### 2. 用户表 `users`

```sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(32) NOT NULL UNIQUE COMMENT '用户名',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    role ENUM('user', 'admin') NOT NULL DEFAULT 'user' COMMENT '用户角色',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    
    INDEX idx_username (username),
    INDEX idx_role (role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

---

## 六、安全与并发控制

### 1. JWT 鉴权机制

- 使用 `HS256` 算法生成 Token
- 设置 24 小时过期时间
- 请求头 `Authorization: Bearer <token>` 校验
- 中间件自动解析用户身份并注入上下文

### 2. 密码安全

- 使用 `bcrypt` 加密存储密码
- 成本因子：`DefaultCost`
- 敏感字段（如 `password_hash`）不返回前端

### 3. 乐观锁并发控制

在多管理员并发更新同一图书时，为避免数据丢失或覆盖问题，系统采用**数据库乐观锁**策略。

#### 实现原理

1. 读取图书信息时携带当前 `version` 值
2. 更新请求中包含原始 `version`
3. 执行更新时使用条件更新语句

该机制在保证数据一致性的同时，避免了悲观锁带来的性能损耗。

---

## TODO 待办事项
- [x] **管理员功能开发**
    - 增删改查、批量查询
- [x] **用户与管理员认证模块**
    - 用户/管理员注册
    - 登录鉴权（JWT 或 Session）
- [ ] **用户端功能开发**
    - 图书借阅功能
    - 图书归还功能
    - 借阅记录表设计与持久化
- [ ] **其他优化**
    - 分页查询支持（`/books/list?page=1&size=10`）
    - 日志系统接入（zap / logrus）
    - 全文检索优化（如集成 Elasticsearch）
    - ...

---

## TODO 代码规范性问题
- [ ] 依赖耦合问题
   - 问题：Controller直接使用config.DB全局变量，难以测试
   - 解决：引入ServiceFactory和依赖注入容器
- [ ] Service层缺少抽象
   - 问题：Service层没有接口，难以Mock测试
   - 解决：为每个Service定义接口，便于单元测试
- [ ] 错误处理不统一
   - 问题：错误处理逻辑分散，错误信息硬编码
   - 解决：定义统一的ServiceError类型和错误处理
- [ ] 测试难度大
   - 问题：大量集成测试，单元测试困难
   - 解决：通过接口抽象实现Mock测试
- [ ] 职责混乱
   - 问题：Controller层包含业务验证逻辑
   - 解决：将验证逻辑下沉到Service层