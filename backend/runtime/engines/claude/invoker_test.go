package claude

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractResultFromLogPrefersResultEvent(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "claude.jsonl")
	content := `{"type":"assistant","message":{"content":[{"type":"text","text":"draft"}]}}` + "\n" +
		`{"type":"result","result":"final","is_error":false}` + "\n"
	if err := os.WriteFile(logPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write log: %v", err)
	}

	if got := ExtractResultFromLog(logPath); got != "final" {
		t.Fatalf("got %q, want final", got)
	}
}

func TestExtractResultFromLogFallsBackToAssistantText(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "claude.jsonl")
	content := `{"type":"assistant","message":{"content":[{"type":"text","text":"answer"}]}}` + "\n"
	if err := os.WriteFile(logPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write log: %v", err)
	}

	if got := ExtractResultFromLog(logPath); got != "answer" {
		t.Fatalf("got %q, want answer", got)
	}
}
