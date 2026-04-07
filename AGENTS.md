# SingerOS Agent Development Guidelines

This document contains essential information for AI agents working with the SingerOS codebase.

## BUILD/LINT/TEST COMMANDS

### Build Commands
- `go build -o ./bundles/singer ./backend/cmd/singer/main.go` - Build the main SingerOS backend service (output to `./bundles/`)
- `go build -o ./bundles/skill-proxy ./backend/cmd/skill-proxy/main.go` - Build the Skill Proxy service (output to `./bundles/`)
- `make docker-build` - Build Docker image (tag: registry.yygu.cn/insmtx/SingerOS:latest)
- `make docker-run` - Run the Docker image locally
- `make run` - Start docker-compose services in foreground mode
- `make run-detached` - Start docker-compose services in detached mode (background)
- `make stop` - Stop docker-compose services
- `make logs` - View docker-compose service logs

### Test Commands
- `go test ./...` - Run all tests in the project
- `go test -v ./...` - Run all tests with verbose output
- `go test ./backend/path/to/package` - Run tests for a specific package
- `go test -run ^TestFunctionName$ ./backend/path` - Run a specific test function
- `go test -race ./...` - Run all tests with race condition detection
- `go test -cover ./...` - Run tests and display coverage information

### Lint Commands
- `go fmt ./...` - Format all Go code
- `go vet ./...` - Vet all Go code for common mistakes
- `golint ./...` - Lint all Go code (install via `go install golang.org/x/lint/golint@latest`)
- `gofmt -s -w .` - Simplify code and write changes (as per the existing Makefile)
- `staticcheck ./...` - Comprehensive Go static analysis (if installed)

## CODE STYLE GUIDELINES

### Import Organization
- Group imports with blank lines between standard library, third-party, and project-specific packages
- Use semantic import aliases only when they prevent naming conflicts
- Organize in three groups: stdlib, third-party, internal packages
```
import (
	"fmt"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	
	"github.com/insmtx/SingerOS/backend/config"
)
```

### Formatting Conventions
- Use tabs for indentation, not spaces (as verified from existing Go files)
- Execute `go fmt ./...` before committing
- Keep lines under 120 characters where possible
- Use `gofmt -s` for simplification of code

### Naming Conventions
- Use CamelCase for exported functions/types (`GetUser`, `UserService`)
- Use camelCase for unexported/internal functions/types (`getUser`, `userService`)
- Use clear, descriptive names; prefer clarity over brevity
- Use consistent names for similar concepts across packages
- Variables related to the system should reference SingerOS concepts

### Types and Interfaces
- Define interfaces close to their first usage
- Keep interfaces small, typically one or a few methods
- Name interface types with "-er" suffix when applicable (e.g., `Runner`, `Handler`)
- Use concrete types explicitly in function signatures when interface is not needed
- Prefer returning pointers for structs when passing to functions if they will be modified

### Error Handling
- Handle errors explicitly; don't ignore them
- Use specific error types when appropriate with wrapped errors
- Follow the pattern: "if err != nil { return err }"
- Use `errors.New()` for simple static strings
- Use `fmt.Errorf()` with `%w` verb for wrapping errors with more context
- Log errors contextually when appropriate

### Additional Guidelines
- All public functions must have GoDoc comments
- Comments should be in English and explain why rather than what
- Maintain consistent logging format throughout the application
- Use context.Context appropriately for cancellation and request-scoped values
- Follow dependency injection patterns rather than global variables
- Use Cobra for command-line interface implementations as shown in main.go files

### Commit Guidelines
- Follow conventional commits format: `<type>(<scope>): <subject>`
- Use Chinese for commit messages in SingerOS project
- Type options include:
  - `feat`: New feature
  - `fix`: Bug fixes
  - `docs`: Documentation updates
  - `style`: Code style adjustments
  - `refactor`: Code refactoring
  - `test`: Testing related
  - `chore`: Build tool or auxiliary tool changes
- When applicable, include detailed descriptions in the body covering technical implementation and business logic

## PROJECT STRUCTURE

- `/backend` - Main Go application code
  - `/backend/cmd/singer` - Main SingerOS backend service entry point
  - `/backend/cmd/skill-proxy` - Skill Proxy service entry point
  - `/backend/config` - Configuration loading and types
  - `/backend/gateway` - HTTP gateway package
  - `/backend/interaction` - Event-driven interaction layer
    - `/backend/interaction/connectors` - Channel connectors (GitHub implemented; GitLab, WeWork stubs)
    - `/backend/interaction/eventbus` - Event bus abstraction (RabbitMQ implementation)
    - `/backend/interaction/gateway` - Event gateway setup
  - `/backend/skills` - Skill interface, types, and examples
  - `/backend/types` - Core domain types (DigitalAssistant, Event, etc.)
- `/bundles` - Build output directory (generated; gitignored)
- `/deployments/build/Dockerfile` - Container build configuration
- `/docs` - Documentation files
- `/proto` - Protobuf definitions
- `/gen` - Generated protobuf Go/Node code
- `/frontend` - Frontend application

## CONTRIBUTION NOTES

- See CONTRIBUTING.md for commit message style guidance
- Make sure all tests pass (`go test ./...`) before submitting changes
- Follow Go's idiomatic patterns and standard practices
- When implementing, consider how components fit into the broader architecture described in ARCHITECTURE.md

## Core Components and Architecture

Based on the AI OS architecture described in ARCHITECTURE.md, the SingerOS platform consists of the following primary components:

1. **Event Gateway** - Receives external events from various channels (✅ implemented)
2. **Event Bus** - Message queue system for decoupling components (✅ RabbitMQ implemented)
3. **Orchestrator** - Core scheduling and coordination mechanism (🔄 planned)
4. **DigitalAssistant** - Top-level abstraction representing AI workers (✅ types defined)
5. **Agent** - Decision-making entities within DigitalAssistants (🔄 planned)
6. **Skill** - Reusable capabilities that can be invoked (✅ interface and base implementation done)
7. **Skill Proxy** - Isolated skill execution service (✅ service skeleton implemented)
8. **Model Router** - Multi-provider LLM routing (🔄 planned)
9. **Memory System** - Short-term and long-term memory (🔄 planned)

## Skill System Definition

Skills represent core building blocks in SingerOS. The `Skill` interface is defined in `backend/skills/skill.go`:

```go
type Skill interface {
    Info() *SkillInfo
    Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)
    Validate(input map[string]interface{}) error
    GetID() string
    GetName() string
    GetDescription() string
}
```

`SkillInfo` contains the skill's metadata:

```
skill.id
skill.name
skill.description
skill.version
skill.category
skill.skill_type       // local | remote
skill.input_schema
skill.output_schema
skill.permissions
```

Embed `BaseSkill` to reduce boilerplate when implementing a new skill.

### Skill Categories

- **Integration Skills** - External system integrations (GitHub, GitLab, WeChat, Feishu, Jira)
- **AI Skills** - LLM-based reasoning capabilities (code_review, summarize, classification)
- **Tool Skills** - Utility capabilities (run_shell, execute_python, http_request)
- **Workflow Skills** - Complex coordinated operations (pr_review_workflow, bug_triage_workflow)

## CHANNEL INTEGRATION

Support for multiple interaction channels via the `Connector` interface in `backend/interaction/connectors/connector.go`:

- **GitHub** (✅ implemented) - Webhook, event parsing, signature verification
- **GitLab** (🔄 stub)
- **Enterprise WeChat / WeWork** (🔄 stub)
- **Feishu** (🔄 planned)
- **App / Webhook** (🔄 planned)

Each channel implements the `Connector` interface:

```go
type Connector interface {
    ChannelCode() string
    RegisterRoutes(r gin.IRouter)
}
```

Events are normalized into the `interaction.Event` type and published to the Event Bus (RabbitMQ).

## Permissions and Security

Granular permissions control at multiple levels:

- DigitalAssistant
- Agent
- Skill
- Tool

Permission model: RBAC + Capability

## GOLANG ENGINE STRUCTURE

Actual code structure as of the current implementation:

```
SingerOS/
│
├── backend/
│   ├── cmd/
│   │   ├── singer/          # Main backend service (HTTP + event gateway)
│   │   └── skill-proxy/     # Skill Proxy service
│   │
│   ├── config/              # Config loading and types (GitHub app config, etc.)
│   │
│   ├── gateway/             # HTTP gateway (placeholder for future routes)
│   │
│   ├── interaction/         # Event-driven interaction layer
│   │   ├── connectors/
│   │   │   ├── github/      # GitHub webhook connector (✅ implemented)
│   │   │   ├── gitlab/      # GitLab connector (🔄 stub)
│   │   │   └── wework/      # WeWork/企业微信 connector (🔄 stub)
│   │   ├── eventbus/
│   │   │   └── rabbitmq/    # RabbitMQ publisher (✅ implemented)
│   │   └── gateway/         # Event gateway router setup
│   │
│   ├── skills/              # Skill interface, BaseSkill, SkillManager interface
│   │   └── examples/        # Example skill implementation
│   │
│   └── types/               # Core domain types
│       ├── digital_assistant.go          # DigitalAssistant, AssistantConfig
│       ├── digital_assistant_instance.go # DigitalAssistantInstance
│       ├── event.go                      # Event (persisted)
│       └── tables.go                    # DB table name constants
│
├── proto/                   # Protobuf definitions
├── gen/                     # Generated code from protos
├── frontend/                # Frontend application
├── deployments/             # Docker build configs
└── docs/                    # Documentation
```

## MINIMUM VISION PRODUCT (MVP)

The initial MVP focuses on these key components:

1. Event Gateway (✅ done)
2. Event Bus / RabbitMQ (✅ done)
3. Skill System interface (✅ done)
4. GitHub Integration (✅ webhook + event parsing done)
5. Skill Proxy service (✅ service skeleton done)
6. Orchestrator (🔄 planned)
7. Agent Engine (🔄 planned)
8. CodeAssistantDigitalAssistant (🔄 planned)

MVP Features:

- PR automatic Review (🔄 planned)
- PR automatic summary (🔄 planned)
- Issue automatic reply (🔄 planned - GitHub issue_comment event supported)
- Code explanation (🔄 planned)

## Technical Stack

Current and planned stack:

| Component | Technology | Status |
|-----------|-----------|--------|
| Language | Golang | ✅ Active |
| HTTP Framework | Gin | ✅ Active |
| CLI Framework | Cobra | ✅ Active |
| Message Queue | RabbitMQ | ✅ Active |
| ORM | GORM | ✅ Active (types defined) |
| Database | Postgres | 🔄 Planned |
| Cache | Redis | 🔄 Planned |
| Vector Store | Qdrant | 🔄 Planned |
| LLM | OpenAI / Claude / DeepSeek | 🔄 Planned |

## DEVELOPMENT WORKFLOW

This section outlines the standard development workflow for contributing to the SingerOS project.

### Standard Development Process

When developing new features or fixing issues, follow these steps:

1. **Synchronize with upstream repository** to ensure you have the latest changes before starting development
2. **Create feature branch** based on the updated main branch
3. **Commit and push changes** to your personal forked repository
4. **Submit pull request manually** via the GitHub web interface to the main repository

### Detailed Steps

#### 1. Synchronize with Upstream Repository

First, make sure your fork is up-to-date with the upstream repository:

```bash
# Add upstream repository if not already added
git remote add upstream https://github.com/insmtx/SingerOS.git

# Fetch the latest changes from upstream
git fetch upstream

# Switch to the main branch
git checkout main

# Merge the upstream changes
git merge upstream/main

# Push the updated main branch to your fork
git push origin main
```

#### 2. Create Feature Branch

Create a new branch based on the updated main branch for your development:

```bash
# Create and switch to a new feature branch
git checkout -b feature/descriptive-feature-name

# Or for bug fixes
git checkout -b fix/descriptive-fix-description

# For more specific features or enhancements
git checkout -b feat/scope-descriptive-name
```

Follow the naming convention for branches:
- `feature/` or `feat/` for major features
- `fix/` for bug fixes
- `enhancement/` for improvements to existing functionality
- `docs/` for documentation changes
- `refactor/` for code restructuring without functionality changes

#### 3. Develop and Commit Changes

After completing your development work, commit your changes to your personal repository:

```bash
# Stage your changes
git add .

# Commit with a properly formatted message following conventional commits:
git commit -m "type(scope): concise description of changes"

# Push your feature branch to your personal fork
git push origin feature/descriptive-feature-name
```

Commit message format follows the conventional commits specification:
```
type(scope): description

[optional body]

[optional footer(s)]
```

Common types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect meaning (white-space, formatting, missing semi-colons, etc.)
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `test`: Adding missing tests or correcting existing tests
- `chore`: Other changes that don't modify src or test files

#### 4. Submit Pull Request

Once your changes are pushed to your personal repository:

1. Navigate to the original SingerOS repository on GitHub
2. Click on "Pull Requests" tab
3. Click on "New Pull Request"
4. Select "Compare across forks"
5. Choose your fork and feature branch as the compared branch
6. Verify the changes shown in the diff
7. Fill in the PR title and description following the project's template
8. Submit the pull request
