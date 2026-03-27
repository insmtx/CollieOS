# Skills 系统

SingerOS 的技能系统提供了一种模块化的方式来实现不同的能力和功能。这个系统允许通过标准化接口集成各种功能。

## 组织结构

- `backend/skills/` - 技能系统的核心接口定义和基础组件
  - `skill.go` - 定义技能的基本接口
  - `manager.go` - 技能管理器的实现
  - `type_converter.go` - 类型转换工具
  - `examples/` - 技能实现示例
  - `tool_skills/` - 工具类技能集合
    - `echo_skill/` - 回声技能（示例），返回输入的内容加上 "Echo: " 前缀

## 如何添加新技能

1. 确定技能所属类别（工具类、AI类、集成类等）
2. 在相应类别目录下创建新技能实现
3. 确保新技能实现 `Skill` 接口定义的方法
4. 在主应用初始化时通过 SkillManager 注册新技能

## 技能架构

技能系统遵循以下架构模式：
DigitalAssistant -> Agent -> Skill -> Tool

其中每个技能都实现了统一的接口，允许编排器根据事件需求灵活调用适当的技能。