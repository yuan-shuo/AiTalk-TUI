package tui

import (
	"aitalk/core"
	"aitalk/utils/json"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Update 处理所有消息和事件
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// 处理窗口大小变化
		m.width = msg.Width
		m.height = msg.Height
		
		// 初始化或调整viewport
		m.viewport = viewport.New(msg.Width, msg.Height-8)
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()
		
		// 调整textarea
		m.textarea.SetWidth(msg.Width - 6)
		m.textarea.SetHeight(3)

	case tea.KeyMsg:
		// 处理键盘事件
		switch m.mode {
		case ModeNormal:
			return m.handleNormalMode(msg)
		case ModeInsert:
			return m.handleInsertMode(msg)
		}

	case responseMsg:
		// 处理AI回复
		m.loading = false
		aiResp := string(msg)
		
		// 添加AI消息到列表
		m.messages = append(m.messages, json.Message{
			Role:    "assistant",
			Content: aiResp,
		})
		
		// 更新viewport内容并滚动到底部
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()
		
		// 继续等待下一次回复
		cmds = append(cmds, waitForResponse(m.respCh))

	case errMsg:
		// 处理错误
		m.loading = false
		m.err = msg
		cmds = append(cmds, waitForError(m.errCh))
	}

	// 更新子组件
	if m.mode == ModeInsert {
		newTextarea, cmd := m.textarea.Update(msg)
		m.textarea = newTextarea
		cmds = append(cmds, cmd)
	} else {
		newViewport, cmd := m.viewport.Update(msg)
		m.viewport = newViewport
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// handleNormalMode 处理普通模式的键盘事件
func (m Model) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "i":
		// 进入输入模式
		m.mode = ModeInsert
		m.textarea.Focus()
		return m, textarea.Blink

	case "q", "ctrl+c":
		// 退出程序
		return m, tea.Quit

	case "j", "down":
		// 向下滚动
		m.viewport.LineDown(1)
		return m, nil

	case "k", "up":
		// 向上滚动
		m.viewport.LineUp(1)
		return m, nil

	case "g":
		// 滚动到顶部
		m.viewport.GotoTop()
		return m, nil

	case "G":
		// 滚动到底部
		m.viewport.GotoBottom()
		return m, nil

	case "ctrl+d":
		// 向下翻半页
		m.viewport.HalfViewDown()
		return m, nil

	case "ctrl+u":
		// 向上翻半页
		m.viewport.HalfViewUp()
		return m, nil

	case "ctrl+f", "pgdown":
		// 向下翻页
		m.viewport.PageDown()
		return m, nil

	case "ctrl+b", "pgup":
		// 向上翻页
		m.viewport.PageUp()
		return m, nil
	}

	return m, nil
}

// handleInsertMode 处理输入模式的键盘事件
func (m Model) handleInsertMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg.Type {
	case tea.KeyEsc:
		// 返回普通模式
		m.mode = ModeNormal
		m.textarea.Blur()
		return m, nil

	case tea.KeyEnter:
		// 检查是否按下了Alt键
		if msg.Alt {
			// Alt+Enter 插入换行，传递给textarea处理
			newTextarea, cmd := m.textarea.Update(msg)
			m.textarea = newTextarea
			return m, cmd
		}

		// Enter: 发送消息
		content := strings.TrimSpace(m.textarea.Value())
		if content == "" {
			return m, nil
		}

		// 添加用户消息
		m.messages = append(m.messages, json.Message{
			Role:    "user",
			Content: content,
		})

		// 清空输入框
		m.textarea.Reset()

		// 更新viewport
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()

		// 发送消息到AI（异步）
		m.loading = true
		go m.sendToAI(content)

		return m, tea.Batch(
			waitForResponse(m.respCh),
			waitForError(m.errCh),
		)

	case tea.KeyCtrlC:
		// Ctrl+C 退出
		return m, tea.Quit
	}

	// 其他所有按键都传递给textarea处理
	newTextarea, cmd := m.textarea.Update(msg)
	m.textarea = newTextarea
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// sendToAI 异步发送消息到AI
func (m *Model) sendToAI(content string) {
	// 准备请求
	if m.setting != "" {
		// 保存原始消息
		originalMessages := m.chatReq.Messages

		// 创建新的消息列表，将角色设定作为第一条
		newMessages := []json.Message{
			{Role: "system", Content: m.setting},
		}

		// 添加原始消息（跳过系统消息）
		for i, msg := range originalMessages {
			if i > 0 {
				newMessages = append(newMessages, msg)
			}
		}

		// 如果是第一次对话，添加开场白
		if m.isFirst && m.prologue != "" {
			newMessages = append(newMessages, json.Message{Role: "assistant", Content: m.prologue})
		}

		// 更新请求消息
		m.chatReq.Messages = newMessages
	}

	// 构建存档文件完整路径
	arcFilePath := m.arcDir
	if !strings.HasSuffix(arcFilePath, "/") && !strings.HasSuffix(arcFilePath, "\\") {
		arcFilePath += "/"
	}
	arcFilePath += m.arcFile

	// 调用核心对话函数
	aiResp, err := core.Chat(content, m.chatReq, m.config, arcFilePath, m.isFirst, m.prologue)
	if err != nil {
		m.errCh <- err
		return
	}

	// 标记第一次对话结束
	m.isFirst = false

	// 发送回复到通道
	m.respCh <- aiResp
}
