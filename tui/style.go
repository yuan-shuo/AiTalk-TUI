package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// 颜色定义
var (
	// 主色调
	colorPrimary   = lipgloss.Color("#7C3AED") // 紫色
	colorSecondary = lipgloss.Color("#10B981") // 绿色
	colorAccent    = lipgloss.Color("#F59E0B") // 橙色
	
	// 文本颜色
	colorText      = lipgloss.Color("#E5E7EB") // 浅灰
	colorTextDim   = lipgloss.Color("#9CA3AF") // 深灰
	colorTextMuted = lipgloss.Color("#6B7280") // 更深灰
	
	// 背景色
	colorBg        = lipgloss.Color("#1F2937") // 深灰背景
	colorBgDarker  = lipgloss.Color("#111827") // 更深背景
	
	// 消息气泡颜色
	colorUserMsg   = lipgloss.Color("#3B82F6") // 蓝色 - 用户消息
	colorAgentMsg  = lipgloss.Color("#10B981") // 绿色 - AI消息
	colorSystemMsg = lipgloss.Color("#6B7280") // 灰色 - 系统消息
)

// 样式定义
var (
	// 标题样式
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorPrimary).
		Padding(0, 1)

	// 副标题样式
	subtitleStyle = lipgloss.NewStyle().
		Foreground(colorTextDim).
		Padding(0, 1)

	// 状态栏样式
	statusBarStyle = lipgloss.NewStyle().
		Background(colorBgDarker).
		Foreground(colorText).
		Padding(0, 1)

	// 模式指示器样式
	modeNormalStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#1F2937")).
		Background(lipgloss.Color("#60A5FA")).
		Padding(0, 1)

	modeInsertStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#1F2937")).
		Background(lipgloss.Color("#34D399")).
		Padding(0, 1)

	// 消息样式
	userNameStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorUserMsg)

	agentNameStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorAgentMsg)

	systemNameStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorSystemMsg)

	// 消息内容样式
	messageContentStyle = lipgloss.NewStyle().
		Foreground(colorText)

	// 用户消息气泡
	userBubbleStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#1E3A5F")).
		Foreground(colorText).
		Padding(0, 1).
		MarginLeft(4)

	// AI消息气泡
	agentBubbleStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#064E3B")).
		Foreground(colorText).
		Padding(0, 1)

	// 输入框样式
	inputBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(0, 1)

	// 输入框焦点样式
	inputBoxFocusedStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorSecondary).
		Padding(0, 1)

	// 帮助文本样式
	helpStyle = lipgloss.NewStyle().
		Foreground(colorTextMuted).
		Italic(true)

	// 错误样式
	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#EF4444")).
		Bold(true)

	// 加载中样式
	loadingStyle = lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true)

	// 分隔线样式
	dividerStyle = lipgloss.NewStyle().
		Foreground(colorTextMuted)
)

// getModeStyle 根据当前模式返回对应的样式
func getModeStyle(mode Mode) lipgloss.Style {
	if mode == ModeInsert {
		return modeInsertStyle
	}
	return modeNormalStyle
}

// getModeText 根据当前模式返回模式文本
func getModeText(mode Mode) string {
	if mode == ModeInsert {
		return " INSERT "
	}
	return " NORMAL "
}

// getHelpText 根据当前模式返回帮助文本
func getHelpText(mode Mode) string {
	if mode == ModeInsert {
		return "Enter发送 | Shift+Enter换行 | Esc返回普通模式"
	}
	return "i-输入模式 | j/↓-上滚 | k/↑-下滚 | q-退出"
}
