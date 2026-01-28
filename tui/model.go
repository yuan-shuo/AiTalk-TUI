package tui

import (
	"aitalk/config"
	"aitalk/utils/json"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Mode 定义TUI的操作模式
type Mode int

const (
	ModeNormal Mode = iota // 普通模式 - 可以滚动浏览
	ModeInsert             // 输入模式 - 可以输入消息
)

// Model 是TUI的主数据模型
type Model struct {
	// UI组件
	viewport viewport.Model // 对话历史显示区域
	textarea textarea.Model // 输入框

	// 数据
	messages   []json.Message // 对话消息列表
	roleName   string         // 角色名
	playerName string         // 玩家名
	arcFile    string         // 存档文件名（不含路径）
	arcDir     string         // 存档目录

	// 配置和依赖
	config   *config.Config // 配置
	chatReq  *json.ChatReq  // 对话请求对象
	rolePath string         // 角色目录路径
	setting  string         // 角色设定
	prologue string         // 角色开场白
	isFirst  bool           // 是否是第一次对话

	// 状态
	mode    Mode  // 当前模式
	err     error // 错误信息
	loading bool  // 是否正在等待AI回复
	width   int   // 终端宽度
	height  int   // 终端高度

	// 内容更新通道
	sendCh chan string // 发送消息的通道
	respCh chan string // 接收AI回复的通道
	errCh  chan error  // 错误通道
}

// 等待AI回复的命令
func waitForResponse(ch chan string) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-ch)
	}
}

// 等待错误的命令
func waitForError(ch chan error) tea.Cmd {
	return func() tea.Msg {
		return errMsg(<-ch)
	}
}

// 消息类型定义
type responseMsg string
type errMsg error

// NewModel 创建新的TUI模型
func NewModel(
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
) *Model {
	return &Model{
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
		sendCh:     make(chan string),
		respCh:     make(chan string),
		errCh:      make(chan error),
	}
}
