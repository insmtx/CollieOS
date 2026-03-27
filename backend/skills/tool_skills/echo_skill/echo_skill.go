// echoskill 包提供回声技能的实现
//
// Echo 技能用于简单返回输入值，但在前部加上 "Echo: " 前缀。
package tool_skills

import (
	"context"
	"fmt"
	"github.com/insmtx/SingerOS/backend/skills"
)

// EchoSkill 是一个回声技能，它会返回输入值但加上 "Echo: " 前缀
type EchoSkill struct {
	skills.BaseSkill
}

// NewEchoSkill 创建一个新的 Echo 技能实例
func NewEchoSkill() *EchoSkill {
	return &EchoSkill{
		BaseSkill: skills.BaseSkill{
			InfoData: &skills.SkillInfo{
				ID:          "echo.simple_echo",
				Name:        "Simple Echo Skill",
				Description: "回声技能，返回相同的输入内容加上 'Echo: ' 前缀",
				Version:     "1.0.0",
				Category:    "tool",
				Author:      "SingerOS Team",
				SkillType:   skills.LocalSkill,
				Permissions: []skills.Permission{},
				InputSchema: skills.InputSchema{
					Type:     "object",
					Required: []string{"input"},
					Properties: map[string]*skills.Property{
						"input": {
							Type:        "string",
							Title:       "输入内容",
							Description: "要回声的内容",
						},
					},
				},
				OutputSchema: skills.OutputSchema{
					Type:     "object",
					Required: []string{"output"},
					Properties: map[string]*skills.Property{
						"output": {
							Type:        "string",
							Title:       "输出内容",
							Description: "经过Echo处理后的内容",
						},
					},
				},
			},
		},
	}
}

// Execute 执行 Echo 技能的业务逻辑
func (s *EchoSkill) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	// 获取必要的参数
	inputVal, exists := input["input"]
	if !exists {
		return nil, fmt.Errorf("必要参数 'input' 未提供")
	}

	inputStr, ok := inputVal.(string)
	if !ok {
		return nil, fmt.Errorf("'input' 参数必须为字符串类型")
	}

	// 执行回声逻辑：加上 "Echo: " 前缀
	output := fmt.Sprintf("Echo: %s", inputStr)

	result := map[string]interface{}{
		"output": output,
	}

	return result, nil
}
