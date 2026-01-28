package chat

import (
	"aitalk/config"
	"aitalk/core"
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

func loopScan(req *json.ChatReq, c *config.Config, arcFilePath string, rolePath string, roleBaseName string) error {

	// 修改角色卡名称
	agentPrompt = fmt.Sprintf("%s:", roleBaseName)
	// 修改玩家名称
	userPrompt = fmt.Sprintf("%s:", c.Player.Name)

	in := bufio.NewScanner(os.Stdin)

	// 读取角色的 setting 文件内容
	settingContent := readRoleSetting(rolePath, roleBaseName)

	// 读取角色的 prologue 文件内容
	prologueContent := readRolePrologue(rolePath, roleBaseName)

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

// 读取角色的 setting 文件内容
func readRoleSetting(rolePath string, roleBaseName string) string {
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

// 读取角色的 prologue 文件内容
func readRolePrologue(rolePath string, roleBaseName string) string {
	if roleBaseName == "" {
		return ""
	}

	prologuePath := filepath.Join(rolePath, roleBaseName+".role", "prologue")
	content, err := os.ReadFile(prologuePath)
	if err != nil {
		return ""
	}

	return string(content)
}
