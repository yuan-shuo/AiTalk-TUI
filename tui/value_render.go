package tui

import (
	"strings"
)

// 统一维护所有变量定义
type VarMap map[string]string

// 变量渲染表统一维护源
func (m Model) GetVars() VarMap {
	return VarMap{
		// 代表将 {user} 替换为玩家名称
		"user": m.playerName,
		// 代表将 {role} 替换为角色名称
		"role": m.roleName,
	}
}

// 替换器只关心 map，不关心具体字段
// 使用 text := ReplaceTextVarWithModelValues("你好 {user}，我是 {role}", model.GetVars())
func ReplaceTextVarWithModelValues(text string, vars VarMap) string {
	result := text
	for key, val := range vars {
		result = strings.ReplaceAll(result, "{"+key+"}", val)
	}
	return result
}

// 以下是旧版本不易维护，暂时弃用

// func (m Model) GetReplaceContext() ReplaceContext {
// 	return ReplaceContext{
// 		PlayerName: m.playerName,
// 		RoleName:   m.roleName,
// 	}
// }

// // ReplaceWithContext 使用上下文替换变量（便于扩展）
// // 可以传入额外的临时变量，不会影响到替换器本身的变量表
// type ReplaceContext struct {
// 	PlayerName string
// 	RoleName   string
// 	// 后续可以添加更多上下文信息
// }

// // ReplaceWithContext 使用上下文进行变量替换
// func ReplaceWithContext(text string, ctx ReplaceContext) string {
// 	result := text

// 	// 替换玩家名称
// 	if ctx.PlayerName != "" {
// 		result = strings.ReplaceAll(result, "{user}", ctx.PlayerName)
// 	}

// 	// 替换角色名称
// 	if ctx.RoleName != "" {
// 		result = strings.ReplaceAll(result, "{role}", ctx.RoleName)
// 	}

// 	return result
// }

// 以下内容不知道是做什么的，先不管

// // VariableReplacer 变量替换器
// type VariableReplacer struct {
// 	// 变量映射表：key是变量名（如"{user}"），value是对应的值
// 	variables map[string]string
// }

// // NewVariableReplacer 创建新的变量替换器
// func NewVariableReplacer(playerName string) *VariableReplacer {
// 	return &VariableReplacer{
// 		variables: map[string]string{
// 			"{user}": playerName, // 玩家名称
// 			// 后续可以在这里添加更多变量，例如：
// 			// "{role}": roleName,      // 角色名称
// 			// "{time}": currentTime,   // 当前时间
// 			// "{date}": currentDate,   // 当前日期
// 		},
// 	}
// }

// // AddVariable 添加自定义变量
// func (vr *VariableReplacer) AddVariable(key, value string) {
// 	// 确保变量名以 { 开头，以 } 结尾
// 	if !strings.HasPrefix(key, "{") {
// 		key = "{" + key
// 	}
// 	if !strings.HasSuffix(key, "}") {
// 		key = key + "}"
// 	}
// 	vr.variables[key] = value
// }

// // Replace 替换文本中的所有变量
// func (vr *VariableReplacer) Replace(text string) string {
// 	result := text
// 	for key, value := range vr.variables {
// 		result = strings.ReplaceAll(result, key, value)
// 	}
// 	return result
// }
