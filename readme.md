# 简易图书管理系统

## 一、功能需求

提供 RESTful API 管理图书信息，使用go-gin框架完成服务端搭建，使用json序列化协议
* 能够将图书数据保存到 MySQL 数据库
* 能够检索图书列表、查看图书详细信息
* 包含完整的单元测试

> 不包含前端 UI 界面。
```
├── internal  ## 内部逻辑
│   ├── api
│   │   ├── req_resp.go
│   │   └── result
│   │       ├── code.go
│   │       └── result.go
│   ├── config
│   │   └── config.go
│   ├── es
│   │   └── client.go
│   ├── handler
│   │   ├── bookHandler.go
│   │   ├── book_controller_test.go
│   │   ├── userHandler.go
│   │   └── user_controller_test.go
│   ├── middleware
│   │   └── auth.go
│   ├── model
│   │   ├── book.go
│   │   └── user.go
│   ├── repo
│   │   └── dao
│   │       ├── bookDAO.go
│   │       ├── book_dao_test.go
│   │       ├── db.go
│   │       ├── userDAO.go
│   │       └── user_dao_test.go
│   ├── router
│   │   └── router.go
│   ├── service
│   │   ├── bookESService.go
│   │   ├── bookService.go
│   │   └── userService.go
│   └── utils
│       ├── jwt.go
│       └── jwt_test.go
├── main.go
└── config.yaml
```

---

## 二、功能模块

| 模块 | 功能 |
|------|------|
| 图书管理 | 增删改查、批量查询、乐观锁控制 |
| 用户系统 | 用户/管理员注册、登录、角色权限控制 |
| 安全认证 | JWT 鉴权、密码加密存储、接口权限校验 |
| 可测试性 | 核心逻辑覆盖单元测试 |
| 部署与构建 |Dockerfile 容器化、Makefile 一键构建部署|

---

## TODO 待办事项
- [x] **管理员功能开发**
    - 增删改查、批量查询
- [x] **用户与管理员认证模块**
    - 用户/管理员注册
    - 登录鉴权（JWT 或 Session）
- [x] **其他优化**
    - 分页查询支持（`/books/list?page=1&size=10`）
    - 全文检索优化（如集成 Elasticsearch）
    - Dockerfile、Makefile
    - ...

---

## TODO 代码规范性问题
- [x] 依赖耦合问题
   - 问题：Controller直接使用config.DB全局变量，难以测试
   - 解决：引入ServiceFactory和依赖注入容器
- [x] Service层缺少抽象
   - 问题：Service层没有接口，难以Mock测试
   - 解决：为每个Service定义接口，便于单元测试
- [x] 错误处理不统一
   - 问题：错误处理逻辑分散，错误信息硬编码
   - 解决：定义统一的ServiceError类型和错误处理
- [x] 测试难度大
   - 问题：大量集成测试，单元测试困难
   - 解决：通过接口抽象实现Mock测试
- [x] 职责混乱
   - 问题：Controller层包含业务验证逻辑
   - 解决：将验证逻辑下沉到Service层