package tui

import (
	"aitalk/utils/json"
	"fmt"
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
	subtitle := subtitleStyle.Render(fmt.Sprintf("📁 %s", m.arcFile))
	
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

// renderUserMessage 渲染用户消息
func (m Model) renderUserMessage(content string) string {
	name := userNameStyle.Render(fmt.Sprintf("[%s]", m.playerName))
	
	// 处理多行内容
	lines := strings.Split(content, "\n")
	var contentLines []string
	for _, line := range lines {
		contentLines = append(contentLines, messageContentStyle.Render(line))
	}
	renderedContent := strings.Join(contentLines, "\n")

	// 用户消息右对齐
	bubble := userBubbleStyle.Render(renderedContent)
	
	// 计算缩进使消息右对齐
	availableWidth := m.width - lipgloss.Width(bubble) - 2
	if availableWidth < 0 {
		availableWidth = 0
	}
	
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Right, 
		lipgloss.JoinVertical(lipgloss.Right, name, bubble))
}

// renderAgentMessage 渲染AI消息
func (m Model) renderAgentMessage(content string) string {
	name := agentNameStyle.Render(fmt.Sprintf("[%s]", m.roleName))
	
	// 处理多行内容
	lines := strings.Split(content, "\n")
	var contentLines []string
	for _, line := range lines {
		contentLines = append(contentLines, messageContentStyle.Render(line))
	}
	renderedContent := strings.Join(contentLines, "\n")

	bubble := agentBubbleStyle.Render(renderedContent)
	
	return lipgloss.JoinVertical(lipgloss.Left, name, bubble)
}

// renderSystemMessage 渲染系统消息
func (m Model) renderSystemMessage(content string) string {
	name := systemNameStyle.Render("[System]")
	
	// 处理多行内容
	lines := strings.Split(content, "\n")
	var contentLines []string
	for _, line := range lines {
		contentLines = append(contentLines, messageContentStyle.Render(line))
	}
	renderedContent := strings.Join(contentLines, "\n")

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
	
	// 使用PlaceHorizontal来布局
	statusContent := lipgloss.JoinHorizontal(lipgloss.Left, left, 
		strings.Repeat(" ", m.width-lipgloss.Width(left)-lipgloss.Width(right)-4), right)
	
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
