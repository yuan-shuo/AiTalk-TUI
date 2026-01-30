package tui

import (
	"strings"
)

// VariableReplacer 变量替换器
type VariableReplacer struct {
	// 变量映射表：key是变量名（如"{user}"），value是对应的值
	variables map[string]string
}

// NewVariableReplacer 创建新的变量替换器
func NewVariableReplacer(playerName string) *VariableReplacer {
	return &VariableReplacer{
		variables: map[string]string{
			"{user}": playerName, // 玩家名称
			// 后续可以在这里添加更多变量，例如：
			// "{role}": roleName,      // 角色名称
			// "{time}": currentTime,   // 当前时间
			// "{date}": currentDate,   // 当前日期
		},
	}
}

// AddVariable 添加自定义变量
func (vr *VariableReplacer) AddVariable(key, value string) {
	// 确保变量名以 { 开头，以 } 结尾
	if !strings.HasPrefix(key, "{") {
		key = "{" + key
	}
	if !strings.HasSuffix(key, "}") {
		key = key + "}"
	}
	vr.variables[key] = value
}

// Replace 替换文本中的所有变量
func (vr *VariableReplacer) Replace(text string) string {
	result := text
	for key, value := range vr.variables {
		result = strings.ReplaceAll(result, key, value)
	}
	return result
}

// ReplaceWithContext 使用上下文替换变量（便于扩展）
// 可以传入额外的临时变量，不会影响到替换器本身的变量表
type ReplaceContext struct {
	PlayerName string
	RoleName   string
	// 后续可以添加更多上下文信息
}

// ReplaceWithContext 使用上下文进行变量替换
func ReplaceWithContext(text string, ctx ReplaceContext) string {
	result := text

	// 替换玩家名称
	if ctx.PlayerName != "" {
		result = strings.ReplaceAll(result, "{user}", ctx.PlayerName)
	}

	// 替换角色名称
	if ctx.RoleName != "" {
		result = strings.ReplaceAll(result, "{role}", ctx.RoleName)
	}

	return result
}
