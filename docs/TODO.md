# 🚀 SingerOS 后端开发 TODO（融合最终版）

> 目标：2周内完成 GitHub PR 自动 Review MVP
> 原则：**先闭环，再抽象；先能跑，再优雅**

---

# 🧭 总体阶段划分

| 阶段      | 目标          | 时间     |
| ------- | ----------- | ------ |
| Phase 1 | 跑通最小闭环（必须）  | Week 1 |
| Phase 2 | 抽象核心能力（可扩展） | Week 2 |
| Phase 3 | MVP业务完善     | Week 2 |
| Phase 4 | 可扩展能力（延后）   | 可选     |

---

# 🔴 Phase 1：最小闭环（必须完成）

> 核心目标：**PR → 自动评论（端到端跑通）**

---

## 1.1 基础运行环境

### Task 1: Docker化项目（优先级最高）

**目标**

* 一键启动所有服务（singer + rabbitmq）

**验收**

* `docker-compose up` 能启动
* singer服务可访问
* RabbitMQ正常

---

## 1.2 GitHub Connector（只做最小）

### Task 2: GitHub Webhook接入

**接口**

```
POST /webhook/github
```

**实现**

* 验证签名
* 解析 PR opened 事件
* 转换为统一 Event

**Event结构（统一标准）**

```go
type Event struct {
    ID        string
    Type      string   // "pr_opened"
    Source    string   // "github"
    Payload   map[string]interface{}
}
```

**验收**

* 能打印PR事件日志
* 能正确解析 repo / pr_number

---

## 1.3 EventBus（RabbitMQ）

### Task 3: 事件发布

**实现**

* GitHub Connector → RabbitMQ

**Topic**

```
github.pr.opened
```

---

### Task 4: 事件消费（Orchestrator初版）

**实现**

* 从RabbitMQ消费事件
* 打日志验证

**验收**

* 收到PR事件

---

## 1.4 LLM能力（必须抽象！）

### Task 5: LLM Provider（简化版）

```go
type LLM interface {
    Generate(ctx context.Context, prompt string) (string, error)
}
```

实现：

* OpenAI Provider（先写死）

**验收**

* 能调用LLM返回结果

---

## 1.5 第一个Skill（核心）

### Task 6: PR Review Skill（本地实现）

**输入**

```json
{
  "diff": "...",
}
```

**输出**

```json
{
  "comments": ["xxx", "xxx"]
}
```

**逻辑**

* prompt + LLM

---

## 1.6 GitHub API能力

### Task 7: 评论PR

**能力**

* Create PR Comment

---

## 1.7 Orchestrator（最简版）

### Task 8: 事件 → Skill（硬编码）

```go
if event.Type == "pr_opened" {
    runPRReview()
}
```

流程：

```
Event → 获取PR diff → 调用Review Skill → 发评论
```

---

## ✅ Phase 1验收标准（必须达成）

✔ 提交PR → 自动生成评论
✔ 全链路日志可见
✔ 无panic

---

# 🟡 Phase 2：核心抽象（避免后期推翻）

> 这一阶段只做**必要抽象，不做过度设计**

---

## 2.1 Skill体系标准化（关键！）

### Task 9: 定义Skill接口

```go
type Skill interface {
    Name() string
    Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)
}
```

---

### Task 10: Skill Manager（核心）

```go
type SkillManager struct {
    registry map[string]Skill
}
```

能力：

* Register
* Get
* Execute

---

## 2.2 Orchestrator升级

### Task 11: 去硬编码

改为：

```go
type Handler func(ctx context.Context, event Event)

map[event.Type]Handler
```

---

## 2.3 Agent（轻量版，不做复杂）

### Task 12: 引入Agent（但简化）

```go
type Agent interface {
    Run(ctx context.Context, event Event) error
}
```

实现：

* PRReviewAgent

👉 注意：

* ❌ 不做Plan
* ❌ 不做Reflect

---

# 🟢 Phase 3：MVP业务完善

---

## 3.1 GitHub能力补全

### Task 13: 获取PR Diff

---

### Task 14: 支持Issue Comment

---

## 3.2 第二个Skill

### Task 15: Issue Reply Skill

---

## 3.3 CodeAssistant

### Task 16: Digital Assistant

职责：

* 绑定：

  * PR → PRReviewAgent
  * Issue → IssueReplyAgent

---

## 3.4 完整Workflow

### Task 17: PR Review流程

```
Webhook → Event → Orchestrator → Agent → Skill → GitHub API
```

---

## ✅ Phase 3验收

✔ PR自动Review
✔ Issue自动回复
✔ 多事件支持

---

# 🔵 Phase 4：扩展能力（不要提前做）

---

## 4.1 Skill Proxy（远程化）

👉 只有当你需要：

* 多实例
* 多语言Skill

再做

---

## 4.2 Memory系统

👉 只有当：

* 多轮对话
* 长上下文

再做

---

## 4.3 多Agent编排

👉 等复杂任务再说

---

# 📁 推荐目录结构（已优化）

```
backend/
├── cmd/
│   └── singer/
├── interaction/
│   └── connectors/github/
├── eventbus/
├── orchestrator/
├── agent/
│   └── pr_review/
├── skills/
│   ├── manager.go
│   ├── github/
│   └── ai/
├── llm/
├── integration/github/
```

---

# 🧪 强制开发规范（必须执行）

每个Task必须：

* [ ] `go build ./...`
* [ ] `go test ./...`
* [ ] 有最小测试
* [ ] 日志可观测

---

# 🧨 风险控制（重点）

### ❗ 不允许做的事情

* ❌ 不要一开始搞 Skill Proxy
* ❌ 不要设计复杂 Agent（Plan/Reflect）
* ❌ 不要做 Memory
* ❌ 不要做多租户
* ❌ 不要做抽象过度的 Workflow Engine

---

# 🎯 最终MVP定义（非常关键）

当满足：

```
1. PR opened
2. 自动分析diff
3. 自动评论
4. Issue自动回复
```

👉 **项目就算成功**

---

# 💡 最后一句（给你团队用的）

> 这个阶段的目标不是“做一个AI平台”，
> 而是**证明这个系统真的能帮人写代码评论**

---

如果你下一步要推进，我可以帮你再做一版：

👉 **“团队分工版本（2-3人怎么拆任务）”**
👉 **“每个Task对应PR粒度（直接能用Git管理）”**
