package builtin

import (
	"testing"

	"github.com/insmtx/SingerOS/backend/config"
	"github.com/insmtx/SingerOS/backend/runtime/engines"
)

func TestNewRegistryFromConfigRegistersEnabledEngines(t *testing.T) {
	registry, err := NewRegistryFromConfig(&config.CLIEnginesConfig{
		Engines: map[string]config.CLIEngineConfig{
			engines.EngineClaude: {Enabled: true, Path: "claude"},
			engines.EngineCodex:  {Enabled: false, Path: "codex"},
		},
	})
	if err != nil {
		t.Fatalf("build registry: %v", err)
	}

	if _, ok := registry.Get(engines.EngineClaude); !ok {
		t.Fatal("expected claude engine to be registered")
	}
	if _, ok := registry.Get(engines.EngineCodex); ok {
		t.Fatal("disabled codex engine should not be registered")
	}
}

func TestNewRegistryFromConfigRejectsUnsupportedEngine(t *testing.T) {
	_, err := NewRegistryFromConfig(&config.CLIEnginesConfig{
		Engines: map[string]config.CLIEngineConfig{
			"unknown": {Enabled: true},
		},
	})
	if err == nil {
		t.Fatal("expected unsupported engine error")
	}
}
