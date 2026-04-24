package engines

import (
	"os"
	"strings"
)

// BuildBaseEnv 返回当前环境并附加额外的环境变量。
func BuildBaseEnv(extraEnv map[string]string) []string {
	env := os.Environ()
	for k, v := range extraEnv {
		if strings.TrimSpace(k) != "" && v != "" {
			env = append(env, k+"="+v)
		}
	}
	return env[:len(env):len(env)]
}

// BuildRunEnv 为 CLI 进程组装环境变量条目，根据不同的模型提供商设置对应的 API 密钥环境变量。
func BuildRunEnv(baseEnv []string, extraEnv []string, model ModelConfig) []string {
	env := make([]string, 0, len(baseEnv)+len(extraEnv)+4)
	env = append(env, baseEnv...)
	env = append(env, extraEnv...)

	switch strings.ToLower(model.Provider) {
	case "anthropic", "claude":
		appendEnvIfSet(&env, "ANTHROPIC_API_KEY", model.APIKey)
		appendEnvIfSet(&env, "ANTHROPIC_AUTH_TOKEN", model.APIKey)
		appendEnvIfSet(&env, "ANTHROPIC_BASE_URL", model.BaseURL)
	case "openai", "codex", "deepseek", "moonshot", "qwen", "zhipu", "":
		appendEnvIfSet(&env, "OPENAI_API_KEY", model.APIKey)
		appendEnvIfSet(&env, "OPENAI_API_BASE", model.BaseURL)
		appendEnvIfSet(&env, "OPENAI_BASE_URL", model.BaseURL)
	default:
		appendEnvIfSet(&env, "OPENAI_API_KEY", model.APIKey)
		appendEnvIfSet(&env, "OPENAI_BASE_URL", model.BaseURL)
	}
	return env
}

func appendEnvIfSet(env *[]string, key string, value string) {
	if value != "" {
		*env = append(*env, key+"="+value)
	}
}
