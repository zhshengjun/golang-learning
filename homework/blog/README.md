# Blog

Go 博客后端学习项目，基于 Gin、GORM 和 MySQL，实现用户注册、登录及用户信息管理。

## 运行环境

- Go 1.26 或更高版本（版本要求见 `go.mod`）
- MySQL 数据库
- 可访问数据库的本地或开发环境

## 安装依赖

必须从项目目录执行命令，因为配置文件使用相对路径 `config.yaml`：

```bash
cd homework/blog
go mod download
```

项目依赖由 `go.mod` 和 `go.sum` 管理，主要包括：

- Gin：HTTP 服务和路由
- GORM：数据库访问
- MySQL Driver：MySQL 驱动
- Viper：读取 YAML 配置
- JWT：登录认证

## 配置

编辑 `config.yaml`，填写数据库连接信息和 JWT 密钥：

```yaml
database:
  dsn: 用户名:密码@tcp(数据库地址:3306)/数据库名?charset=utf8mb4&parseTime=True&loc=Local

jwt:
  secret: 自定义JWT密钥
```

启动前请确认：

1. MySQL 服务已启动。
2. DSN 中的数据库已经创建。
3. 数据库账号具备访问该数据库的权限。
4. 项目当前未执行自动建表，所需数据表需要提前准备。

## 启动

在 `homework/blog` 目录执行：

```bash
go run .
```

服务默认监听 `8080` 端口：

```text
http://localhost:8080
```

也可以编译后启动：

```bash
go build -o blog .
./blog
```

## 接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/user/register` | 用户注册 |
| POST | `/login` | 用户登录 |
| POST | `/logout` | 退出登录，需要登录 |
| GET | `/user/info?id=1` | 查询用户信息，需要登录 |
| POST | `/user/update` | 更新用户信息，需要登录 |
| DELETE | `/user/delete` | 删除用户，需要登录 |

## 当前代码状态

当前工作区执行 `go test ./...` 时，`router/article_router.go` 与 `router/user_router.go` 存在同名函数重复定义，项目会在编译阶段失败。修复重复定义后，才能正常执行上述启动命令。
