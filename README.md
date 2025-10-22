# OJPlus - 在线评测系统

## 项目介绍
OJPlus 是一个现代化的在线评测系统，支持多语言代码提交、自动评测、实时反馈等功能。系统采用微服务架构设计，使用 Go 作为后端主要开发语言，Vue.js 作为前端框架，PostgreSQL 和 Redis 作为数据存储，Kafka 作为消息队列以支持多判题机管理。

## 系统架构

### 后端架构 (Go)
- **API 服务**：处理前端请求，提供 RESTful API
- **判题服务**：负责代码编译、执行和结果评估
- **Kafka 消息队列**：实现判题任务的分发和负载均衡
- **权限管理服务**：用户认证与授权
- **数据统计服务**：记录和分析用户提交和系统性能

### 前端架构 (Vue.js)
- **用户界面**：题目浏览、代码编辑器、提交历史等
- **管理界面**：题目管理、用户管理、系统配置等
- **实时反馈**：显示提交状态和判题结果

### 数据存储
- **PostgreSQL**：存储用户信息、题目数据、提交记录等
- **Redis**：缓存热点数据、会话管理、排行榜等

### 消息队列
- **Kafka**：连接 API 服务和判题服务，实现异步处理和多判题机管理

## 功能特性
- 支持多种编程语言的代码评测
- 实时评测结果反馈
- 灵活的测试用例配置
- 详细的提交历史和错误分析
- 用户排行榜和竞赛管理
- 可水平扩展的判题系统

## 技术栈
- 后端：Go, Gin, GORM
- 前端：Vue.js, Vuetify
- 数据库：PostgreSQL, Redis
- 消息队列：Kafka
- 容器化：Docker, Kubernetes

## 开发与部署

### 环境要求
- Go 1.18+
- Node.js 14+
- PostgreSQL 13+
- Redis 6+
- Kafka 2.8+
- Docker & Docker Compose (可选)

### 开发设置
```bash
# 克隆代码库
git clone https://github.com/yourusername/ojplus.git
cd ojplus

# 后端设置
cd backend
go mod download
go run main.go

# 前端设置
cd ../frontend
npm install
npm run serve
```

### 部署
```bash
# 使用 Docker Compose
docker-compose up -d
```

## 贡献指南
欢迎提交问题和功能请求！请遵循以下步骤：

1. Fork 该仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个 Pull Request

## 许可证
[MIT License](LICENSE)