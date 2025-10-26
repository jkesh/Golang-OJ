# Postman 测试指南

本指南将帮助你使用 Postman 测试 OJPlus 系统的代码提交功能。

## 1. 用户认证

在提交代码之前，你需要先进行用户认证。

### 1.1 用户登录

**请求方法**: POST  
**URL**: `http://localhost:8080/api/auth/login`  
**Headers**: 
- Content-Type: application/json

**Body** (JSON):
```json
{
  "username": "testuser",
  "password": "testpassword"
}
```

**响应示例**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "Test User"
  }
}
```

保存返回的 token，后续请求需要在 headers 中添加 Authorization。

### 1.2 注册新用户（如果需要）

**请求方法**: POST  
**URL**: `http://localhost:8080/api/auth/register`  
**Headers**: 
- Content-Type: application/json

**Body** (JSON):
```json
{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "newpassword",
  "nickname": "New User"
}
```

## 2. 代码提交

### 2.1 提交代码

**请求方法**: POST  
**URL**: `http://localhost:8080/api/submissions`  
**Headers**: 
- Content-Type: application/json
- Authorization: Bearer <your_token_here>

**Body** (JSON):
```json
{
  "problem_id": 1,
  "language": "go",
  "code": "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}"
}
```

**参数说明**:
- `problem_id`: 题目ID（需要确保题目存在）
- `language`: 编程语言（支持: go, cpp, java, python）
- `code`: 提交的代码内容

**响应示例**:
```json
{
  "message": "提交代码",
  "status": "success"
}
```

## 3. 查看提交记录

### 3.1 获取所有提交记录

**请求方法**: GET  
**URL**: `http://localhost:8080/api/submissions`  
**Headers**: 
- Authorization: Bearer <your_token_here>

**响应示例**:
```json
{
  "submissions": [
    {
      "id": 1,
      "problem_id": 1,
      "user_id": 1,
      "language": "go",
      "code": "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}",
      "status": "pending",
      "run_time": 0,
      "memory": 0,
      "submitted_at": "2025-10-26T10:00:00Z",
      "judged_at": "0001-01-01T00:00:00Z",
      "error_message": ""
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 10,
  "message": "获取所有提交列表",
  "status": "success"
}
```

### 3.2 获取特定提交详情

**请求方法**: GET  
**URL**: `http://localhost:8080/api/submissions/1`  
**Headers**: 
- Authorization: Bearer <your_token_here>

**响应示例**:
```json
{
  "message": "获取提交详情",
  "id": "1",
  "status": "success"
}
```

## 4. 题目相关接口

### 4.1 获取题目列表

**请求方法**: GET  
**URL**: `http://localhost:8080/api/problems`  
**Headers**: 
- Authorization: Bearer <your_token_here>

### 4.2 获取特定题目详情

**请求方法**: GET  
**URL**: `http://localhost:8080/api/problems/1`  
**Headers**: 
- Authorization: Bearer <your_token_here>

## 5. 注意事项

1. 确保后端服务正在运行，并且端口正确（默认是 8080）
2. 确保 Kafka 服务正在运行，地址为 `111.228.56.17:9092`
3. 确保 PostgreSQL 数据库正在运行
4. 确保 Redis 服务正在运行
5. 提交代码前需要先登录获取 token
6. 在 headers 中正确添加 Authorization 头，格式为 `Bearer <token>`
7. 确保提交的题目 ID 在数据库中存在