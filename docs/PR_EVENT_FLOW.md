# GitHub PR 事件处理流程验证

## 闭环验证清单

### ✓ Webhook 接收层 
- [x] `/github/webhook` 端点已注册
- [x] 支持多种 GitHub 事件类型 (包括 pull_request)
- [x] 签名验证机制有效

### ✓ 事件解析层
- [x] 支持 PR 事件 (`pull_request`) 类型解析
- [x] 自动映射为内部 `interaction.Event` 结构
- [x] 正确填充 Payload 数据

### ✓ 事件发布层
- [x] 发布到特定主题 (`interaction.github.pull_request`)
- [x] 经 RabbitMQ 消息队列传输

### ✓ 事件消费层 (新增)
- [x] Orchestrator 消费 PR 事件主题
- [x] 日志记录与基础验证功能
- [x] 预留扩展处理逻辑的空间

## 主要代码结构变更

### 1. 新增 Orchestrator 组件 (`/backend/orchestrator/`)
- `orchestrator.go`: 核心事件消费者
- `orchestrator_test.go`: 单元测试
- `README.md`: 功能说明文档

### 2. 扩展 GitHub 事件处理 (`/backend/interaction/connectors/github/`)
- `events.go`: 支持 `pull_request` 事件类型转换
- `types.go`: 定义 `EventTypePullRequest` 常量
- `webhook.go`: 根据事件类型发布到不同主题

### 3. 事件主题常量 (`/backend/interaction/topic.go`)
- 新增 `TopicGithubPullRequest` 常量定义

### 4. 主程序集成 (`/backend/cmd/singer/main.go`)
- 集成 Orchestrator 启动流程
- 在服务启动时启动事件流消费

## 完整性测试
以下命令可验证系统构建和运行：
```bash
go build -o ./bundles/singer ./backend/cmd/singer/main.go
go test ./backend/orchestrator/...
```

## 验证示例
当 GitHub 发送 `pull_request` 事件到 `/github/webhook` 时，系统日志将显示：
```
Processing GitHub pull request event: ...
```

## 总结
通过添加 Orchestrator 组件，SingerOS 的 GitHub 集成已形成闭环：从 webhook 接收 → 事件标准化 → 消息发布 → 事件消费 → 日志验证。