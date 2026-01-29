package chat

import (
	"aitalk/config"
	"aitalk/core"
	"aitalk/tui"
	"aitalk/utils/json"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

var (
	userPrompt  string = `YOU (enter "Ctrl+C" to exit):`
	agentPrompt string = `AGENT:`
)

// useTUI 控制是否使用TUI界面，默认启用
var useTUI = true

func loopScan(req *json.ChatReq, c *config.Config, arcFilePath string, rolePath string, roleID string, roleName string) error {
	// 读取角色的 setting 文件内容（通过roleID定位）
	settingContent := readRoleSetting(rolePath, roleID)

	// 读取角色的 prologue 文件内容（通过roleID定位）
	prologueContent := readRolePrologue(rolePath, roleID)

	// 如果使用TUI模式
	if useTUI {
		return runTUI(req, c, arcFilePath, rolePath, roleID, roleName, settingContent, prologueContent)
	}

	// 否则使用传统命令行模式
	return runSimpleMode(req, c, arcFilePath, rolePath, roleID, roleName, settingContent, prologueContent)
}

// runTUI 运行TUI模式
func runTUI(req *json.ChatReq, c *config.Config, arcFilePath string, rolePath string, roleID string, roleName string, settingContent string, prologueContent string) error {
	// 获取存档目录和文件名
	arcDir := filepath.Dir(arcFilePath)
	arcFile := filepath.Base(arcFilePath)

	// 运行TUI（使用roleName作为显示名）
	return tui.Run(
		req.Messages,
		roleName, // 使用显示名
		c.Player.Name,
		arcFile,
		arcDir,
		c,
		req,
		rolePath,
		settingContent,
		prologueContent,
	)
}

// runSimpleMode 运行传统简单模式
func runSimpleMode(req *json.ChatReq, c *config.Config, arcFilePath string, rolePath string, roleID string, roleName string, settingContent string, prologueContent string) error {
	// 修改角色卡名称（使用显示名）
	agentPrompt = fmt.Sprintf("%s:", roleName)
	// 修改玩家名称
	userPrompt = fmt.Sprintf("%s:", c.Player.Name)

	in := bufio.NewScanner(os.Stdin)

	// 打印角色开场白
	if c.Character.Prologue.Enabled && prologueContent != "" {
		fmt.Printf("%s %s\n\n", agentPrompt, prologueContent)
	}

	// 标记是否是第一次对话
	isFirstDialogue := true

	for {
		fmt.Printf("%s ", userPrompt)
		if !in.Scan() {
			break
		}

		text := in.Text()

		// 在每轮对话的 req 变量最前端加入角色设定
		if settingContent != "" {
			// 保存原始消息
			originalMessages := req.Messages

			// 创建新的消息列表，将角色设定作为第一条
			newMessages := []json.Message{
				{Role: "system", Content: settingContent},
			}

			// 添加原始消息（跳过系统消息）
			for i, msg := range originalMessages {
				if i > 0 {
					newMessages = append(newMessages, msg)
				}
			}

			// 如果是第一次对话，添加开场白
			if isFirstDialogue && prologueContent != "" {
				newMessages = append(newMessages, json.Message{Role: "assistant", Content: prologueContent})
			}

			// 更新 req.Messages
			req.Messages = newMessages
		}

		aiResp, err := core.Chat(text, req, c, arcFilePath, isFirstDialogue, prologueContent)
		if err != nil {
			return err
		}
		fmt.Printf("\n%s %s\n\n", agentPrompt, aiResp)

		// 第一次对话结束后，设置为 false
		isFirstDialogue = false
	}
	return nil
}

// 读取角色的 setting 文件内容（通过roleID定位）
func readRoleSetting(rolePath string, roleID string) string {
	if roleID == "" {
		return ""
	}

	settingPath := filepath.Join(rolePath, roleID+".role", "setting")
	content, err := os.ReadFile(settingPath)
	if err != nil {
		return ""
	}

	return string(content)
}

// 读取角色的 prologue 文件内容（通过roleID定位）
func readRolePrologue(rolePath string, roleID string) string {
	if roleID == "" {
		return ""
	}

	prologuePath := filepath.Join(rolePath, roleID+".role", "prologue")
	content, err := os.ReadFile(prologuePath)
	if err != nil {
		return ""
	}

	return string(content)
}
