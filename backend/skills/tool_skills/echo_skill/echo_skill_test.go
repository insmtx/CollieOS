package tool_skills

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEchoSkill_Execute(t *testing.T) {
	skill := NewEchoSkill()

	ctx := context.Background()

	tests := []struct {
		name        string
		input       map[string]interface{}
		expected    map[string]interface{}
		expectError bool
	}{
		{
			name: "valid input string",
			input: map[string]interface{}{
				"input": "Hello World",
			},
			expected: map[string]interface{}{
				"output": "Echo: Hello World",
			},
			expectError: false,
		},
		{
			name: "another valid input",
			input: map[string]interface{}{
				"input": "Testing Echo Skill",
			},
			expected: map[string]interface{}{
				"output": "Echo: Testing Echo Skill",
			},
			expectError: false,
		},
		{
			name:        "missing input parameter",
			input:       map[string]interface{}{},
			expected:    nil,
			expectError: true,
		},
		{
			name: "wrong input type",
			input: map[string]interface{}{
				"input": 123, // not string
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "non-string type in input",
			input: map[string]interface{}{
				"input": []string{"hello", "world"},
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := skill.Execute(ctx, tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestEchoSkill_Info(t *testing.T) {
	skill := NewEchoSkill()

	info := skill.Info()

	assert.Equal(t, "echo.simple_echo", info.ID)
	assert.Equal(t, "Simple Echo Skill", info.Name)
	assert.Equal(t, "回声技能，返回相同的输入内容加上 'Echo: ' 前缀", info.Description)
	assert.Equal(t, "1.0.0", info.Version)
	assert.Equal(t, "tool", info.Category)
	assert.Equal(t, "SingerOS Team", info.Author)
}

func TestEchoSkill_BasicMethods(t *testing.T) {
	skill := NewEchoSkill()

	assert.Equal(t, "echo.simple_echo", skill.GetID())
	assert.Equal(t, "Simple Echo Skill", skill.GetName())
	assert.Equal(t, "回声技能，返回相同的输入内容加上 'Echo: ' 前缀", skill.GetDescription())
}
