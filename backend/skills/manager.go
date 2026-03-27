package skills

import (
	"context"
	"fmt"
	"sync"
)

// SimpleSkillManager 是一个简单的技能管理器实现
type SimpleSkillManager struct {
	skills map[string]Skill
	mutex  sync.RWMutex
}

// NewSimpleSkillManager 创建一个新的 simpleSkillManager 实例
func NewSimpleSkillManager() *SimpleSkillManager {
	return &SimpleSkillManager{
		skills: make(map[string]Skill),
	}
}

// Register 注册一个新技能到管理器
func (sm *SimpleSkillManager) Register(skill Skill) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	skillID := skill.GetID()
	if _, exists := sm.skills[skillID]; exists {
		return fmt.Errorf("skill with ID %s already exists", skillID)
	}

	sm.skills[skillID] = skill
	return nil
}

// Get 根据ID获取已注册的技能
func (sm *SimpleSkillManager) Get(skillID string) (Skill, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	skill, exists := sm.skills[skillID]
	if !exists {
		return nil, fmt.Errorf("skill with ID %s not found", skillID)
	}

	return skill, nil
}

// List 返回所有已注册的技能
func (sm *SimpleSkillManager) List() []Skill {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	skillsList := make([]Skill, 0, len(sm.skills))
	for _, skill := range sm.skills {
		skillsList = append(skillsList, skill)
	}

	return skillsList
}

// Execute 执行指定技能ID的技能
func (sm *SimpleSkillManager) Execute(ctx context.Context, skillID string, input map[string]interface{}) (map[string]interface{}, error) {
	skill, err := sm.Get(skillID)
	if err != nil {
		return nil, err
	}

	if err := skill.Validate(input); err != nil {
		return nil, fmt.Errorf("invalid input provided to skill %s: %w", skillID, err)
	}

	return skill.Execute(ctx, input)
}
