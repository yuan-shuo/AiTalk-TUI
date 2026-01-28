package tui

import (
	"aitalk/config"
	"aitalk/utils/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Run 启动TUI对话界面
func Run(
	messages []json.Message,
	roleName string,
	playerName string,
	arcFile string,
	arcDir string,
	config *config.Config,
	chatReq *json.ChatReq,
	rolePath string,
	setting string,
	prologue string,
) error {
	// 创建初始模型
	m := initialModel(
		messages,
		roleName,
		playerName,
		arcFile,
		arcDir,
		config,
		chatReq,
		rolePath,
		setting,
		prologue,
	)

	// 创建Bubble Tea程序
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // 使用备用屏幕缓冲区
		tea.WithMouseCellMotion(), // 启用鼠标支持
	)

	// 运行程序
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI运行错误: %w", err)
	}

	return nil
}

// initialModel 创建并初始化模型
func initialModel(
	messages []json.Message,
	roleName string,
	playerName string,
	arcFile string,
	arcDir string,
	config *config.Config,
	chatReq *json.ChatReq,
	rolePath string,
	setting string,
	prologue string,
) Model {
	// 初始化textarea
	ta := textarea.New()
	ta.Placeholder = "输入消息..."
	ta.SetHeight(3)
	ta.SetWidth(80)
	ta.ShowLineNumbers = false
	ta.CharLimit = 0 // 无字符限制

	// 配置textarea的快捷键
	ta.KeyMap.InsertNewline.SetEnabled(true)

	// 初始化viewport
	vp := viewport.New(80, 20)

	// 创建模型
	m := Model{
		messages:   messages,
		roleName:   roleName,
		playerName: playerName,
		arcFile:    arcFile,
		arcDir:     arcDir,
		config:     config,
		chatReq:    chatReq,
		rolePath:   rolePath,
		setting:    setting,
		prologue:   prologue,
		isFirst:    true,
		mode:       ModeNormal,
		textarea:   ta,
		viewport:   vp,
		sendCh:     make(chan string),
		respCh:     make(chan string),
		errCh:      make(chan error),
		width:      80,
		height:     24,
	}

	// 设置初始内容
	vp.SetContent(m.renderMessages())
	vp.GotoBottom()

	return m
}

// Init 初始化（实现tea.Model接口）
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		waitForResponse(m.respCh),
		waitForError(m.errCh),
	)
}

// ReadRoleSetting 读取角色设定文件
func ReadRoleSetting(rolePath string, roleBaseName string) string {
	if roleBaseName == "" {
		return ""
	}

	settingPath := filepath.Join(rolePath, roleBaseName+".role", "setting")
	content, err := os.ReadFile(settingPath)
	if err != nil {
		return ""
	}

	return string(content)
}

// ReadRolePrologue 读取角色开场白文件
func ReadRolePrologue(rolePath string, roleBaseName string) string {
	if roleBaseName == "" {
		return ""
	}

	prologuePath := filepath.Join(rolePath, roleBaseName+".role", "prologue")
	content, err := os.ReadFile(prologuePath)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(content))
}
