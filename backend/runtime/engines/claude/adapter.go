// Package claude 将 Claude Code 适配到 SingerOS 外部 CLI 引擎接口。
package claude

import (
	"context"

	"github.com/insmtx/SingerOS/backend/runtime/engines"
)

// Adapter 通过 Claude Code 执行提示。
type Adapter struct {
	invoker *Invoker
}

// NewAdapter 创建 Claude Code 引擎适配器。
func NewAdapter(binary string, extraEnv map[string]string) *Adapter {
	if binary == "" {
		binary = "claude"
	}
	return &Adapter{invoker: NewInvoker(binary, extraEnv)}
}

// Prepare 执行 Claude 工作区设置（当前为空实现）。
func (a *Adapter) Prepare(_ context.Context, _ engines.PrepareRequest) error {
	return nil
}

// Run 启动 Claude Code 并返回进程句柄。
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
