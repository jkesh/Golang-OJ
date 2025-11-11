# OJPlus - 在线评测系统

## 项目介绍

OJPlus 是一个现代化的在线评测系统，支持多语言代码提交、自动评测、实时反馈等功能。系统采用前后端分离架构设计，使用 Go 作为后端主要开发语言，Vue.js 作为前端框架，PostgreSQL 作为主数据库，Redis 作为缓存数据库，Kafka 作为消息队列以支持多判题机管理。

## 系统架构

### 后端架构 (Go)
- **API 服务**：基于 Gin 框架构建，处理前端请求，提供 RESTful API
- **数据库访问**：使用 GORM ORM 库操作 PostgreSQL 数据库
- **认证授权**：基于 JWT 实现用户认证与权限控制
- **消息队列**：集成 Kafka 实现判题任务的异步处理和分发
- **配置管理**：使用 Viper 管理系统配置

### 前端架构 (Vue.js)
- **用户界面**：基于 Element Plus 组件库构建，包含题目浏览、代码编辑、提交记录等功能
- **状态管理**：使用 Vuex 管理应用状态
- **路由管理**：使用 Vue Router 实现前端路由
- **代码编辑器**：集成 Monaco Editor 提供代码编辑功能

### 数据存储
- **PostgreSQL**：存储用户信息、题目数据、提交记录等核心数据
- **Redis**：缓存热点数据、会话管理等

### 消息队列
- **Kafka**：连接 API 服务和判题服务，实现异步处理和多判题机管理

## 功能特性
- 多语言代码评测支持 (Go, C++, Java, Python)
- 实时代码提交与评测结果反馈
- 用户认证与权限管理
- 题目管理与测试用例配置
- 提交记录查看与详细结果分析
- 管理员后台系统
- 多判题机支持与管理

## 技术栈

### 后端技术栈
- Go 1.25+
- Gin Web 框架
- GORM ORM 库
- PostgreSQL 数据库
- Redis 缓存
- Kafka 消息队列
- Viper 配置管理
- JWT 认证

### 前端技术栈
- Vue.js 3.x
- Element Plus UI 组件库
- Vue Router 路由管理
- Vuex 状态管理
- Axios HTTP 客户端
- Monaco Editor 代码编辑器
- Vite 构建工具

## 项目结构

```
.
├── backend                 # 后端代码
│   ├── api                 # API 控制器
│   ├── auth                # 认证相关
│   ├── common              # 公共组件
│   ├── config              # 配置管理
│   ├── middleware          # 中间件
│   ├── models              # 数据模型
│   ├── go.mod             # Go 模块定义
│   └── main.go            # 程序入口
├── frontend                # 前端代码
│   ├── components          # 公共组件
│   ├── src                 # 源代码
│   │   ├── api             # API 客户端
│   │   ├── components      # Vue 组件
│   │   ├── router          # 路由配置
│   │   ├── store           # Vuex 状态管理
│   │   └── views           # 页面视图
│   ├── package.json       # Node.js 依赖配置
│   └── vite.config.js     # Vite 配置
├── POSTMAN_TESTING.md     # Postman 测试指南
└── README.md              # 项目说明文档
```

## 开发环境搭建

### 后端环境要求
- Go 1.25+
- PostgreSQL 13+
- Redis 6+
- Kafka 2.8+

### 前端环境要求
- Node.js 14+
- npm 或 yarn

### 后端开发设置
```bash
# 进入后端目录
cd backend

# 下载依赖
go mod download

# 运行程序
go run main.go
```

### 前端开发设置
```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

## 配置说明

后端配置文件位于 `backend/config/config.json`，包含以下主要配置项：
- 服务器配置 (端口、主机等)
- 数据库配置 (主机、端口、用户名、密码等)
- Redis配置 (主机、端口、密码等)
- Kafka配置 (Brokers、主题等)
- 判题配置 (超时时间、内存限制、支持的语言等)

## API 测试

项目提供了 Postman 测试指南 [POSTMAN_TESTING.md](POSTMAN_TESTING.md)，包含以下测试用例：
- 用户认证接口测试
- 代码提交接口测试
- 提交记录查询接口测试
- 题目相关接口测试

## 部署说明

### Docker 部署 (推荐)

```bash
# 构建并启动所有服务
docker-compose up -d
```

### 手动部署

1. 配置并启动 PostgreSQL、Redis、Kafka 服务
2. 修改配置文件 `backend/config/config.json`
3. 编译并运行后端服务
4. 构建并部署前端应用

## 开发指南

### 添加新功能

1. 后端：
   - 在 `backend/models` 中添加数据模型
   - 在 `backend/api` 中添加 API 控制器
   - 在 `backend/api/api.go` 中注册路由

2. 前端：
   - 在 `frontend/src/api` 中添加 API 客户端
   - 在 `frontend/src/views` 中添加页面组件
   - 在 `frontend/src/router/index.js` 中添加路由

### 代码规范

- 后端遵循 Go 语言编码规范
- 前端遵循 Vue.js 编码规范
- 提交代码前运行相应测试

## 贡献指南

欢迎提交问题和功能请求！请遵循以下步骤：

1. Fork 该仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个 Pull Request

## 许可证

[MIT License](LICENSE)