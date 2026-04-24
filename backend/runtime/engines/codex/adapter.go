// Package codex 将 Codex CLI 适配到 SingerOS 外部 CLI 引擎接口。
package codex

import (
	"context"

	"github.com/insmtx/SingerOS/backend/runtime/engines"
)

// Adapter 通过 Codex CLI 执行提示。
type Adapter struct {
	invoker *Invoker
}

// NewAdapter 创建 Codex CLI 引擎适配器。
func NewAdapter(binary string, extraEnv map[string]string) *Adapter {
	if binary == "" {
		binary = "codex"
	}
	return &Adapter{invoker: NewInvoker(binary, NewSessionStore(), extraEnv)}
}

// Prepare 执行 Codex 工作区设置（当前为空实现）。
func (a *Adapter) Prepare(_ context.Context, _ engines.PrepareRequest) error {
	return nil
}

// Run 启动 Codex CLI 并返回进程句柄。
func (a *Adapter) Run(ctx context.Context, req engines.RunRequest) (*engines.RunHandle, error) {
	proc, events, err := a.invoker.Run(ctx, req)
	if err != nil {
		return nil, err
	}
	return &engines.RunHandle{
		Process:       proc,
		Events:        events,
		ExtractResult: ExtractResultFromLog,
	}, nil
}

var _ engines.Engine = (*Adapter)(nil)
