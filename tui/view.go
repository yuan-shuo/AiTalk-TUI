package tui

import (
	"aitalk/utils/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View 渲染整个界面
func (m Model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	// 构建界面组件
	title := m.renderTitle()
	viewport := m.renderViewport()
	statusBar := m.renderStatusBar()
	inputBox := m.renderInputBox()

	// 组合界面
	var sb strings.Builder
	sb.WriteString(title)
	sb.WriteString("\n")
	sb.WriteString(viewport)
	sb.WriteString("\n")
	sb.WriteString(statusBar)
	sb.WriteString("\n")
	sb.WriteString(inputBox)

	return sb.String()
}

// renderTitle 渲染标题栏
func (m Model) renderTitle() string {
	title := titleStyle.Render(fmt.Sprintf("🎭 %s", m.roleName))
	// // 去掉 .jsonl 后缀显示
	// displayName := strings.TrimSuffix(m.arcFile, filepath.Ext(m.arcFile))
	// subtitle := subtitleStyle.Render(fmt.Sprintf("📁 %s", displayName))

	// 去掉 .jsonl 后缀
	displayName := strings.TrimSuffix(m.arcFile, filepath.Ext(m.arcFile))

	// 去掉 hash 前缀（例如 "a2esd-test" -> "test"）
	if idx := strings.Index(displayName, "-"); idx != -1 {
		displayName = displayName[idx+1:]
	}
	// 顶部显示的对话名
	subtitle := subtitleStyle.Render(fmt.Sprintf("📁 %s", displayName))

	return lipgloss.JoinHorizontal(lipgloss.Left, title, "  ", subtitle)
}

// renderViewport 渲染对话历史区域
func (m Model) renderViewport() string {
	content := m.renderMessages()
	m.viewport.SetContent(content)

	// 设置viewport高度（减去标题、状态栏和输入框的高度）
	viewportHeight := m.height - 8
	if viewportHeight < 5 {
		viewportHeight = 5
	}
	m.viewport.Height = viewportHeight
	m.viewport.Width = m.width

	return m.viewport.View()
}

// renderMessages 渲染所有消息
func (m Model) renderMessages() string {
	var sb strings.Builder

	for _, msg := range m.messages {
		// 跳过system消息（角色设定）
		if msg.Role == "system" {
			continue
		}

		rendered := m.renderMessage(msg)
		sb.WriteString(rendered)
		sb.WriteString("\n\n")
	}

	// 如果正在加载，显示加载提示
	if m.loading {
		sb.WriteString(loadingStyle.Render(fmt.Sprintf("%s 正在思考...", m.roleName)))
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderMessage 渲染单条消息
func (m Model) renderMessage(msg json.Message) string {
	switch msg.Role {
	case "user":
		return m.renderUserMessage(msg.Content)
	case "assistant":
		return m.renderAgentMessage(msg.Content)
	default:
		return m.renderSystemMessage(msg.Content)
	}
}

// wrapText 将文本按指定宽度自动换行
func wrapText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}

	var result []string
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		// 处理每一行
		currentLine := ""
		currentWidth := 0

		// 按字符处理（支持中文）
		for _, char := range line {
			charWidth := 1
			if char > 127 {
				// 中文字符通常占2个宽度
				charWidth = 2
			}

			if currentWidth+charWidth > maxWidth && currentLine != "" {
				result = append(result, currentLine)
				currentLine = string(char)
				currentWidth = charWidth
			} else {
				currentLine += string(char)
				currentWidth += charWidth
			}
		}

		if currentLine != "" {
			result = append(result, currentLine)
		}
	}

	return strings.Join(result, "\n")
}

// escapeLipglossChars 转义 lipgloss 特殊字符
func escapeLipglossChars(text string) string {
	// 转义可能导致渲染问题的字符
	text = strings.ReplaceAll(text, "\x1b", "") // 移除 ESC 字符
	text = strings.ReplaceAll(text, "\x00", "") // 移除空字符
	return text
}

// renderUserMessage 渲染用户消息
func (m Model) renderUserMessage(content string) string {
	name := userNameStyle.Render(fmt.Sprintf("[%s]", m.playerName))

	// 转义特殊字符
	content = escapeLipglossChars(content)

	// 计算可用宽度（考虑边距和气泡样式）
	maxContentWidth := m.width - 10
	if maxContentWidth < 20 {
		maxContentWidth = 20
	}

	// 自动换行处理
	wrappedContent := wrapText(content, maxContentWidth)

	// 直接使用气泡样式渲染，不再逐行处理
	bubble := userBubbleStyle.Render(wrappedContent)

	// 用户消息右对齐
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Right,
		lipgloss.JoinVertical(lipgloss.Right, name, bubble))
}

// renderAgentMessage 渲染AI消息
func (m Model) renderAgentMessage(content string) string {
	name := agentNameStyle.Render(fmt.Sprintf("[%s]", m.roleName))

	// 转义特殊字符
	content = escapeLipglossChars(content)

	// 计算可用宽度（考虑边距和气泡样式）
	maxContentWidth := m.width - 6
	if maxContentWidth < 20 {
		maxContentWidth = 20
	}

	// 自动换行处理
	wrappedContent := wrapText(content, maxContentWidth)

	// 直接使用气泡样式渲染
	bubble := agentBubbleStyle.Render(wrappedContent)

	return lipgloss.JoinVertical(lipgloss.Left, name, bubble)
}

// renderSystemMessage 渲染系统消息
func (m Model) renderSystemMessage(content string) string {
	name := systemNameStyle.Render("[System]")

	// 转义特殊字符
	content = escapeLipglossChars(content)

	// 计算可用宽度
	maxContentWidth := m.width - 4
	if maxContentWidth < 20 {
		maxContentWidth = 20
	}

	// 自动换行处理
	wrappedContent := wrapText(content, maxContentWidth)

	// 直接使用样式渲染
	renderedContent := messageContentStyle.Render(wrappedContent)

	return lipgloss.JoinVertical(lipgloss.Left, name, renderedContent)
}

// renderStatusBar 渲染状态栏
func (m Model) renderStatusBar() string {
	// 模式指示器
	modeIndicator := getModeStyle(m.mode).Render(getModeText(m.mode))

	// 帮助文本
	helpText := helpStyle.Render(getHelpText(m.mode))

	// 消息计数
	msgCount := helpStyle.Render(fmt.Sprintf("%d messages", len(m.messages)))

	// 组合状态栏
	left := lipgloss.JoinHorizontal(lipgloss.Left, modeIndicator, "  ", helpText)
	right := msgCount

	// 计算中间空格数量，确保不为负数
	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	spaceCount := m.width - leftWidth - rightWidth - 4
	if spaceCount < 0 {
		spaceCount = 0
	}

	// 使用PlaceHorizontal来布局
	statusContent := lipgloss.JoinHorizontal(lipgloss.Left, left,
		strings.Repeat(" ", spaceCount), right)

	return statusBarStyle.Width(m.width).Render(statusContent)
}

// renderInputBox 渲染输入框
func (m Model) renderInputBox() string {
	if m.mode == ModeInsert {
		// 输入模式下显示输入框
		inputView := m.textarea.View()
		return inputBoxFocusedStyle.Width(m.width - 4).Render(inputView)
	}

	// 普通模式下显示提示
	hint := helpStyle.Render("按 i 进入输入模式")
	return inputBoxStyle.Width(m.width - 4).Render(hint)
}
