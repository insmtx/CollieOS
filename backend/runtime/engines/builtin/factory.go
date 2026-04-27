// Package builtin 连接内置的外部 CLI 引擎适配器。
package builtin

import (
	"fmt"
	"strings"

	"github.com/insmtx/SingerOS/backend/config"
	"github.com/insmtx/SingerOS/backend/runtime/engines"
	"github.com/insmtx/SingerOS/backend/runtime/engines/claude"
	"github.com/insmtx/SingerOS/backend/runtime/engines/codex"
)

// NewRegistryFromConfig 创建包含所有已启用内置 CLI 引擎的注册表。
func NewRegistryFromConfig(cfg *config.CLIEnginesConfig) (*engines.Registry, error) {
	registry := engines.NewRegistry()
	if cfg == nil {
		return registry, nil
	}
	for name, item := range cfg.Engines {
		if !item.Enabled {
			continue
		}
		engine, err := newEngine(strings.ToLower(strings.TrimSpace(name)), item)
		if err != nil {
			return nil, err
		}
		if err := registry.Register(name, engine); err != nil {
			return nil, err
		}
	}
	return registry, nil
}

func newEngine(name string, cfg config.CLIEngineConfig) (engines.Engine, error) {
	switch name {
	case engines.EngineClaude:
		return claude.NewAdapter(cfg.Path, cfg.ExtraEnv), nil
	case engines.EngineCodex:
		return codex.NewAdapter(cfg.Path, cfg.ExtraEnv), nil
	default:
		return nil, fmt.Errorf("unsupported CLI engine %q", name)
	}
}
