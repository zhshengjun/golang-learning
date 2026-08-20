# Go Blog

Go 博客后端学习项目，基于 Gin、GORM 和 MySQL，实现登录、用户、文章和评论接口。

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

## 登录接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/login` | 用户登录 |
| POST | `/logout` | 退出登录，需要登录 |

登录成功后，服务会通过 `blog_token` Cookie 保存 JWT；除注册和登录外的接口都需要登录。

## 用户管理接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/user/register` | 用户注册 |
| GET | `/user/info?id=1` | 查询用户信息，需要登录 |
| POST | `/user/updated` | 更新用户信息，需要登录 |
| DELETE | `/user/deleted` | 删除用户，需要登录 |

## 文章接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/article/created` | 创建文章 |
| PUT | `/article/updated` | 更新文章 |
| PUT | `/article/published` | 发布文章 |
| DELETE | `/article/deleted` | 删除文章 |
| GET | `/article/info?id=1` | 查询文章详情 |
| GET | `/article/page?currentPage=1&pageSize=10` | 分页查询当前用户的文章 |

## 评论接口

评论接口均需要登录：

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/comment/created` | 创建评论或回复 |
| DELETE | `/comment/deleted` | 删除评论及其回复 |
| GET | `/comment/list?articleId=1` | 查询文章评论列表 |

## 检查项目

在 `homework/blog` 目录执行：

```bash
go test ./...
```

当前项目可以通过编译检查。启动前仍需确保 `config.yaml` 中的 MySQL 服务可访问，并且目标数据库及所需数据表已经准备好；项目启动代码未执行自动建表。
