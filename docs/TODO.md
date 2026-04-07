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

* `docker-compose up` 能启动 ❌ *待完成*
* singer服务可访问
* RabbitMQ正常

---

## 1.2 GitHub Connector（只做最小）

### Task 2: GitHub Webhook接入 ✅ *已完成*

**接口**

```
POST /webhook/github
```

**实现**

* 验证签名 ✅ *已完成*
* 解析 PR opened 事件 ✅ *已完成*
* 转换为统一 Event ✅ *已完成*

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

* 能打印PR事件日志 ✅ *已完成*
* 能正确解析 repo / pr_number ✅ *已完成*

---

## 1.3 EventBus（RabbitMQ）

### Task 3: 事件发布 ✅ *已完成*

**实现**

* GitHub Connector → RabbitMQ ✅ *已完成*

**Topic**

```
github.pr.opened
```

---

### Task 4: 事件消费（Orchestrator初版） ✅ *已完成*

**实现**

* 从RabbitMQ消费事件 ✅ *已完成*
* 打日志验证 ✅ *已完成*

**验收**

* 收到PR事件 ✅ *已完成*

---

## 1.4 LLM能力（必须抽象！）

### Task 5: LLM Provider（简化版） ✅ *已完成*

```go
type LLM interface {
    Generate(ctx context.Context, prompt string) (string, error)
}
```

实现：

* OpenAI Provider（先写死）✅ *已完成*

**验收**

* 能调用LLM返回结果 ✅ *已完成*

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

* Create PR Comment ❌ *待完成*

---

## 1.7 Orchestrator（最简版）

### Task 8: 事件 → Skill（硬编码）✅ *已完成*

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

✔ 提交PR → 自动生成评论 ❌ *待完成* - 还没有实际的PR评论，只有一个echo技能
✔ 全链路日志可见 ✅ *已完成*
✔ 无panic ✅ *已完成*

---

# 🟡 Phase 2：核心抽象（避免后期推翻）

> 这一阶段只做**必要抽象，不做过度设计**

---

## 2.1 Skill体系标准化（关键！）

### Task 9: 定义Skill接口 ✅ *已完成*

```go
type Skill interface {
    Name() string
    Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)
}
```

---

### Task 10: Skill Manager（核心）✅ *已完成*

```go
type SkillManager struct {
    registry map[string]Skill
}
```

能力：

* Register ✅ *已完成*
* Get ✅ *已完成*
* Execute ✅ *已完成*

---

## 2.2 Orchestrator升级

### Task 11: 去硬编码 ❌ *待完成*

改为：

```go
type Handler func(ctx context.Context, event Event)

map[event.Type]Handler
```

*当前仍为硬编码*

---

## 2.3 Agent（轻量版，不做复杂）

### Task 12: 引入Agent（但简化）

```go
type Agent interface {
    Run(ctx context.Context, event Event) error
}
```

实现：

* PRReviewAgent ❌ *待完成*

👉 注意：

* ❌ 不做Plan
* ❌ 不做Reflect

---

# 🟢 Phase 3：MVP业务完善

---

## 3.1 GitHub能力补全

### Task 13: 获取PR Diff ❌ *待完成*

---

### Task 14: 支持Issue Comment ❌ *待完成*

---

## 3.2 第二个Skill

### Task 15: Issue Reply Skill ❌ *待完成*

---

## 3.3 CodeAssistant

### Task 16: Digital Assistant ❌ *待完成*

职责：

* 绑定：

  * PR → PRReviewAgent
  * Issue → IssueReplyAgent

---

## 3.4 完整Workflow

### Task 17: PR Review流程 ❌ *待完成*

```
Webhook → Event → Orchestrator → Agent → Skill → GitHub API
```

---

## ✅ Phase 3验收

✔ PR自动Review ❌ *待完成*
✔ Issue自动回复 ❌ *待完成*
✔ 多事件支持 ❌ *待完成*

---

# 🔵 Phase 4：扩展能力（不要提前做）

---

## 4.1 Skill Proxy（远程化）

👉 只有当你需要：

* 多实例  
* 多语言Skill

再做

✅ *已完成* - Skill Proxy框架已建立

---

## 4.2 Memory系统

👉 只有当：

* 多轮对话
* 长上下文

再做

❌ *未开始*

---

## 4.3 多Agent编排

👉 等复杂任务再说

❌ *未开始*

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

* ✅ `go build ./...` - 已验证可以通过
* ✅ `go test ./...` - 测试框架完整 
* ✅ 有最小测试 - 测试框架存在
* ✅ 日志可观测 - 已集成yg-go/logs

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
1. PR opened ✅ *事件已捕获*
2. 自动分析diff ❌ *待完成 - 需要PR Diff获取和AI分析技能*
3. 自动评论 ❌ *待完成 - 需要AI分析和评论发布技能*
4. Issue自动回复 ❌ *待完成 - 需要AI回复和发布技能*
```

👉 **项目就算成功**

当前状态：基础事件管道已完成，AI业务处理逻辑(MVP核心部分)待完成
---

# 💡 最后一句（给你团队用的）

> 这个阶段的目标不是“做一个AI平台”，
> 而是**证明这个系统真的能帮人写代码评论**

---

# 🔍 当前项目状态摘要 (截至 2026-04-07)

## 已完成的基础设施组件:

### 1. LLM系统 ✅ 
- 接口定义完整 (`backend/llm/provider.go`)
- OpenAI Provider 实现 (`backend/llm/openai/provider.go`) 
- LLM Router 实现，支持多提供商和降级 (`backend/llm/router.go`)

### 2. Skill系统 ✅  
- 完整的Skill接口和BaseSkill抽象 (`backend/skills/skill.go`)
- SkillManager实现，支持注册、获取和执行 (`backend/skills/manager.go`)
- 示例Skills实现 (`backend/skills/examples/`, `backend/skills/tool_skills/`)

### 3. GitHub连接器 ✅
- Webhook接收和验证 (`backend/interaction/connectors/github/webhook.go`)
- 事件解析和转换为统一Event格式 (`backend/interaction/connectors/github/events.go`)
- 支持PR和Issue Comment事件

### 4. 事件总线和编排器 ✅
- 事件发布到RabbitMQ (`backend/interaction/eventbus/`)
- 事件消费和处理 (`backend/orchestrator/orchestrator.go`)
- 预设处理流程，包括PR和Issue事件

### 5. 数据库系统 ✅
- 数据库连接和初始化 (`backend/database/database.go`)
- 自动迁移机制
- 核心模型定义：DigitalAssistant, Event, User等 (`backend/types/`)

### 6. 主服务集成 ✅
- singer主服务 (`backend/cmd/singer/`)
- skill-proxy服务 (`backend/cmd/skill-proxy/`)

## 待完成的核心业务功能:

### 1. PR自动Review功能 ❌
- 获取PR Diff的具体实现
- 代码分析和审查逻辑  
- PR评论发布功能

### 2. Agent引擎 ❌ 
- Agent接口的实现
- 决策和执行逻辑
- 与Skills的编排机制

### 3. 数字助手(DigitalAssistant)配置和管理 ❌
- 数字助手实例的配置和运行
- 绑定PR/Issue事件到具体代理(agent)
- 持久化和管理接口

### 4. 更完整的Skill集成 ❌
- GitHub API调用的Skills（如PR评论、Issue回复）
- AI相关的Skills（代码分析、摘要生成等）

## 总体状态评估:

- Phase 1 基础架构: ✅ **大部分完成** (除Docker化)
- Phase 2 核心抽象: ✅ **Skill系统和LLM完成，去硬编码和Agent待完成**
- Phase 3 MVP业务: ❌ **核心AI业务功能需要补充**
- 风险: 连续性和数据一致性机制需要加强测试

---

如果你下一步要推进，我可以帮你再做一版：

👉 **“团队分工版本（2-3人怎么拆任务）”**
👉 **“每个Task对应PR粒度（直接能用Git管理）”**
