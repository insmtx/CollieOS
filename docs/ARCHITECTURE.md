
# AI OS 架构设计文档

## 1. 项目愿景

本项目旨在构建一个 **企业级 AI 操作系统（AI OS）**，用于管理和运行 **AI Digital Assistants（AI数字员工）**。

核心目标：

> 让企业可以像组织员工一样管理 AI。


## 2. 核心设计原则

### 2.1 事件驱动架构

AI OS 的核心是 **事件驱动系统**：

```
External Systems
    ↓
Event Gateway
    ↓
Event Bus
    ↓
Orchestrator
    ↓
Agent Runtime
    ↓
Skills/Tools
    ↓
Result
```

所有 AI 行为都由 **事件触发**。

### 2.2 DigitalAssistant 是最高抽象

系统核心结构：

```
DigitalAssistant
    ↓
Agent
    ↓
Skill
    ↓
Tool
```

DigitalAssistant 代表一个完整的 AI 数字员工，包含 Agent、Skills、Tools、记忆和知识库等组件。每个 DigitalAssistant 可以配置不同的能力和行为模式。

### 2.3 控制平面 vs 数据平面

SingerOS 严格分离了：

* **控制平面**：管理、注册、存储、策略
* **数据平面**：运行时执行
* **基础设施层**：数据库、消息队列、存储

## 3. 认证与账号管理

### 3.1 OAuth 授权流程

系统支持通过 OAuth 协议集成外部平台（如 GitHub、GitLab）的用户认证。当用户在外部平台完成授权后，系统会获取访问令牌并创建授权账号记录。

授权流程包括：
- 用户发起授权请求
- 重定向到外部平台进行身份验证
- 外部平台返回授权码
- 系统使用授权码换取访问令牌
- 创建或更新授权账号记录

### 3.2 授权账号模型

系统维护授权账号（AuthorizedAccount）的概念，用于存储用户在外部平台的认证信息。每个授权账号关联到具体的渠道类型和用户标识。

授权账号包含平台信息、访问令牌、刷新令牌、过期时间以及关联的内部用户标识。这些信息用于在后续操作中以用户身份执行任务。

### 3.3 账号解析器

在运行时，系统通过账号解析器（AccountResolver）确定执行操作时使用的具体账号。解析器根据事件触发者、数字助手配置或显式指定的账户信息来选择合适的授权账号。

账号解析器支持多种策略，包括使用默认账户、使用事件触发者账户或使用最近使用的账户。

### 3.4 授权选择器

工具执行时通过授权选择器（AuthSelector）确定访问权限和使用的账户。授权选择器综合数字助手配置、Agent 配置、技能权限和工具策略来决定最终的执行身份。

授权选择器确保工具在执行时具备适当的权限，并且符合系统的安全策略要求。

# 4. 核心架构

## 4.1 架构总览

```
                   +----------------------+
                   |   External Systems   |
                   +----------+-----------+
                              |
                              v
                     +----------------+
                     | Event Gateway  |
                     +--------+-------+
                              |
                              v
                      +--------------+
                      | Event Bus    |
                      +------+-------+
                             |
                             v
                     +---------------+
                     | Orchestrator  |
                     +-------+-------+
                             |
           +-----------------+----------------+
           |                                  |
           v                                  v
    +-------------+                   +---------------+
    | Agent       |                   | Skills        |
    | Runtime     |                   | Catalog       |
    +------+------+                   +-------+-------+
           |                                  |
           v                                  v
     +------------+                    +-------------+
     | Tools      |                    | Tool        |
     | Execution  |                    | Registry    |
     +------------+                    +-------------+
```


# 5. 核心组件

## 5.1 Event Gateway

负责接收所有外部事件并标准化为内部事件格式。

统一事件结构包含事件ID、追踪ID、渠道来源、事件类型、触发者、关联仓库、上下文信息和负载数据等字段。


## 5.2 Event Bus

作用：
- 解耦系统
- 支持高并发
- 支持异步


## 5.3 Orchestrator

AI OS 的事件路由和调度器。

职责：

```
事件消费 → 找到匹配的 Handler → 调用 Agent Runtime → 处理结果
```


# 6. Agent Runtime 设计

## 6.1 Runner 接口

定义 Orchestrator 和具体实现之间的抽象边界，提供事件处理能力。

## 6.2 EinoRunner 实现

EinoRunner 是基于 CloudWeGo Eino 框架的 LLM Agent 运行时。

核心组件包含聊天模型、工具适配器、技能上下文、工具上下文和系统提示词。

执行流程：
1. 事件接收
2. 提示构建
3. 技能注入
4. 工具注入
5. LLM 执行
6. 结果处理


## 6.3 DigitalAssistant 设计

DigitalAssistant 是 SingerOS 中数字助手的顶级抽象。

结构包含标识符、组织信息、名称、描述、状态、版本和配置等字段。

AssistantConfig 包含运行时配置、LLM 配置、技能列表、渠道列表、知识库列表、记忆配置和策略配置。



# 7. Skills 系统

Skills 是 SingerOS 的可复用能力单元。

## 7.1 Skill 接口

Skill 接口提供能力信息、执行、验证和标识获取等核心方法。

SkillInfo 结构包含 ID、名称、描述、版本、分类、技能类型、输入输出模式和权限定义。

## 7.2 文件化 Skills

Skills 可以定义为文件系统目录中的 SKILL.md 文件，使用 YAML frontmatter 定义元数据。

Catalog 系统负责扫描和加载这些文件化技能。

## 7.3 Skill 分类

- **Integration Skills** - 外部系统集成
- **AI Skills** - AI 推理能力
- **Tool Skills** - 底层工具能力
- **Workflow Skills** - 组合能力


# 8. Tools 系统

Tools 是 SingerOS 的底层原子能力，提供与外部系统交互的具体实现。

## 8.1 Tool 接口

Tool 接口提供工具信息、执行和验证方法。RuntimeTool 额外支持带执行上下文的工具调用。

ToolInfo 结构包含名称、描述、提供者、只读标志和输入模式定义。

## 8.2 Tool Registry

注册表使用线程安全的映射管理所有工具的注册和查询。

## 8.3 Tool Runtime

运行时负责工具查找、上下文构建、输入验证、执行和结果返回的完整流程。


# 9. Skill Proxy

Skill Proxy 提供独立的技能执行隔离环境。

用途：
- 隔离高风险技能的执行
- 支持多语言技能实现
- 资源限制和沙箱化


# 10. 多交互渠道设计

系统必须支持多渠道交互。

## 10.1 Connector 接口

统一抽象定义渠道代码获取和 HTTP 路由注册方法。


# 11. 记忆系统

## 11.1 短期记忆

短期记忆用于存储会话级别的上下文信息，包括当前对话状态、用户意图和临时变量。短期记忆在会话结束后通常会被丢弃。

短期记忆支持：
- 会话上下文维护
- 多轮对话状态跟踪
- 临时信息缓存

## 11.2 长期记忆

长期记忆用于持久化存储用户偏好、历史经验和学习成果。长期记忆可以跨越多个会话，使 DigitalAssistant 能够记住用户习惯和历史交互。

长期记忆支持：
- 用户偏好存储
- 历史经验回顾
- 知识库更新


# 12. 知识系统

## 12.1 知识库引用

系统支持 DigitalAssistant 关联多个知识库资源。知识库可以包含文档、代码片段、常见问题解答等各种类型的知识内容。

## 12.2 知识访问方式

DigitalAssistant 在执行任务时可以访问关联的知识库，通过检索增强生成（RAG）等技术获取相关知识，提高回答质量和任务执行效果。

知识访问支持：
- 向量检索
- 语义匹配
- 相关度排序


# 13. 配置管理

## 13.1 配置加载

系统通过配置文件加载运行时所需的各项参数。配置按功能模块组织，支持环境变量覆盖。

## 13.2 主要配置项

配置包含以下主要部分：
- **GitHub 配置**：OAuth 客户端信息、Webhook 密钥、应用 ID
- **GitLab 配置**：OAuth 客户端信息、访问令牌
- **RabbitMQ 配置**：连接地址、队列名称、交换器配置
- **数据库配置**：连接字符串、连接池大小
- **LLM 配置**：模型选择、API 密钥、温度参数
- **服务配置**：HTTP 端口、日志级别


# 14. 持久化层

## 14.1 数据库使用

系统使用 GORM 作为 ORM 框架进行数据库操作。GORM 提供了模型定义、迁移支持和查询构建功能。

## 14.2 主要数据表

系统包含以下主要数据表：
- **DigitalAssistant 表**：存储数字助手定义和配置
- **Event 表**：存储事件历史记录
- **AuthorizedAccount 表**：存储授权账号信息
- **Skill 表**：存储技能注册信息
- **Session 表**：存储会话状态和短期记忆


# 15. 事件流示例

典型的从 Webhook 到工具执行的事件流如下：

1. **Webhook 接收**：外部系统（如 GitHub）发送 Webhook 事件到 Event Gateway
2. **事件标准化**：Event Gateway 将外部事件格式转换为系统内部的标准事件格式
3. **事件发布**：事件被发布到 Event Bus（RabbitMQ）的相应队列
4. **事件消费**：Orchestrator 从 Event Bus 消费事件
5. **助手匹配**：Orchestrator 根据事件内容确定需要触发的 DigitalAssistant
6. **Agent 激活**：DigitalAssistant 中的 Agent 被激活，加载相关配置
7. **技能注入**：根据配置将相关的 Skills 和 Tools 注入到 Agent 上下文
8. **LLM 推理**：Agent 调用 LLM 进行推理，决定需要执行的操作
9. **账号解析**：如果需要调用外部工具，通过 AccountResolver 确定使用的授权账号
10. **工具执行**：Agent 调用相应的 Tool 执行具体操作
11. **结果处理**：工具执行结果返回给 Agent，Agent 根据结果决定后续动作
12. **响应返回**：最终结果通过 Event Gateway 返回给调用方或发布到外部系统


# 16. Golang 工程结构

```
SingerOS/
│
├── backend/
│   ├── cmd/
│   │   ├── singer/          # 主服务
│   │   └── skill-proxy/     # Skill Proxy 服务
│   │
│   ├── config/              # 配置加载与类型定义
│   │
│   ├── gateway/             # HTTP gateway
│   │
│   ├── interaction/         # 事件驱动交互层
│   │   ├── connectors/
│   │   ├── eventbus/
│   │   └── gateway/
│   │
│   ├── skills/              # Skill 接口与实现
│   │
│   └── types/               # 核心领域类型
│
├── proto/                   # Protobuf 定义
├── gen/                     # 生成的代码
└── frontend/                # 前端应用
```


# 17. 技术栈

| 组件 | 技术 |
|------|------|
| 语言 | Golang |
| HTTP 框架 | Gin |
| CLI 框架 | Cobra |
| 消息队列 | RabbitMQ |
| ORM | GORM |
| 数据库 | Postgres (规划中) |
| 缓存 | Redis (规划中) |
| 向量库 | Qdrant (规划中) |
| LLM | OpenAI / Claude / DeepSeek (规划中) |


# 18. 未来扩展

系统架构支持未来扩展不同类型的 AI 数字员工，包括但不限于：

- **AI产品经理**：协助需求分析、产品规划和用户研究
- **AI QA工程师**：自动化测试、代码质量检查和问题发现
- **AI DevOps工程师**：持续集成/部署、基础设施管理和监控
- **AI数据分析师**：数据探索、报表生成和洞察发现
- **AI技术支持**：用户问题解答和工单处理

每个数字员工可以配置不同的技能组合、知识领域和交互渠道，满足企业不同场景的需求。


# 总结

本 AI OS 架构的核心：

```
Event Driven
Skill Based
Agent Orchestration
DigitalAssistant
```

系统特点：

```
高度模块化
企业级权限
多渠道交互
可扩展 AI 能力
```

