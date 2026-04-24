package claude

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/insmtx/SingerOS/backend/runtime/engines"
)

// Invoker 启动 Claude Code 进程。
type Invoker struct {
	binary  string
	baseEnv []string
}

// NewInvoker 创建 Claude Code 调用器。
func NewInvoker(binary string, extraEnv map[string]string) *Invoker {
	return &Invoker{
		binary:  binary,
		baseEnv: engines.BuildBaseEnv(extraEnv),
	}
}

type streamEvent struct {
	Type    string         `json:"type"`
	Message *streamMessage `json:"message,omitempty"`
	Result  string         `json:"result,omitempty"`
	IsError bool           `json:"is_error,omitempty"`
}

type streamMessage struct {
	Content []streamContent `json:"content"`
}

type streamContent struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// ExtractResultFromLog 从 stream-json 日志中返回最终的 Claude 结果。
func ExtractResultFromLog(logPath string) string {
	result, _, lastAssistantText := scanResultLog(logPath)
	if result != "" {
		return result
	}
	return lastAssistantText
}

func scanResultLog(logPath string) (result string, isError bool, lastAssistantText string) {
	f, err := os.Open(logPath)
	if err != nil {
		return "", false, ""
	}
	defer f.Close()

	engines.ScanJSONLines(f, func(line string) bool {
		var event streamEvent
		if json.Unmarshal([]byte(line), &event) != nil {
			return true
		}
		switch event.Type {
		case "assistant":
			if event.Message != nil {
				for _, block := range event.Message.Content {
					if block.Type == "text" && block.Text != "" {
						lastAssistantText = block.Text
					}
				}
			}
		case "result":
			result = event.Result
			isError = event.IsError
			return false
		}
		return true
	})
	return result, isError, lastAssistantText
}

// Run 启动 Claude Code 进程并将 stdout/stderr 写入 req.LogPath。
func (inv *Invoker) Run(ctx context.Context, req engines.RunRequest) (engines.Process, <-chan engines.Event, error) {
	if req.LogPath == "" {
		return nil, nil, fmt.Errorf("log path is required")
	}
	args := buildArgs(req)

	logFile, err := os.OpenFile(req.LogPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("open log file: %w", err)
	}

	execCtx := ctx
	cancel := func() {}
	if req.Timeout > 0 {
		execCtx, cancel = context.WithTimeout(ctx, req.Timeout)
	}

	cmd := exec.CommandContext(execCtx, inv.binary, args...)
	cmd.Dir = req.WorkDir
	cmd.Stdin = strings.NewReader(req.Prompt)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Env = engines.BuildRunEnv(inv.baseEnv, req.ExtraEnv, req.Model)

	if err := cmd.Start(); err != nil {
		cancel()
		_ = logFile.Close()
		return nil, nil, fmt.Errorf("start claude: %w", err)
	}

	events := make(chan engines.Event, 2)
	proc := engines.NewCmdProcess(cmd)
	events <- engines.Event{Type: engines.EventStarted}

	go func() {
		defer close(events)
		defer logFile.Close()
		defer cancel()
		if err := cmd.Wait(); err != nil {
			events <- engines.Event{Type: engines.EventError, Content: err.Error()}
			return
		}
		result, isError, _ := scanResultLog(req.LogPath)
		if isError {
			if result == "" {
				result = "claude execution failed"
			}
			events <- engines.Event{Type: engines.EventError, Content: result}
			return
		}
		events <- engines.Event{Type: engines.EventDone}
	}()

	return proc, events, nil
}

func buildArgs(req engines.RunRequest) []string {
	args := []string{
		"--dangerously-skip-permissions",
		"--verbose",
		"--output-format", "stream-json",
	}
	if req.Model.Model != "" {
		args = append(args, "--model", req.Model.Model)
	}
	if req.SessionID != "" {
		if req.Resume {
			args = append(args, "--resume", req.SessionID)
		} else {
			args = append(args, "--session-id", req.SessionID)
		}
	}
	return append(args, "--print")
}
