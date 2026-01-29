package chat

import (
	"aitalk/config"
	aitalkJson "aitalk/utils/json"
	"aitalk/utils/role"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func Run(c *config.Config, arcPath string, rolePath string) error {
	// 确保存档目录存在
	if err := os.MkdirAll(arcPath, 0755); err != nil {
		return err
	}

	// 确保角色目录存在
	if err := os.MkdirAll(rolePath, 0755); err != nil {
		return err
	}

	var req *aitalkJson.ChatReq
	var arcfile string
	var roleID string   // 角色hash ID
	var roleName string // 角色显示名

	// 询问使用什么对话
	arcName, err := askUseWhichChatToStart(arcPath)
	if err != nil {
		return err
	}

	if arcName == newDialogueName {
		// 1. 新建对话逻辑

		// 询问使用什么角色
		selectedRoleID, selectedRoleName, err := askUseWhichRole(rolePath)
		if err != nil {
			return err
		}

		// 处理角色
		if selectedRoleID == "" {
			// 创建新角色
			in := bufio.NewScanner(os.Stdin)
			fmt.Print("请输入新角色的名称: ")
			if !in.Scan() {
				return fmt.Errorf("no input was given")
			}
			newRoleName := in.Text()

			// 调用创建角色函数，返回角色ID
			newRoleID, err := role.CreateRole(rolePath, newRoleName, c.TextEditor)
			if err != nil {
				return err
			}

			roleID = newRoleID
			roleName = newRoleName
		} else {
			roleID = selectedRoleID
			roleName = selectedRoleName
		}

		// 询问对话名称
		in := bufio.NewScanner(os.Stdin)
		fmt.Print("enter the name of this new dialogue: ")
		if !in.Scan() {
			return fmt.Errorf("no input was given")
		}
		dialogueName := in.Text()

		// 拼接存档文件路径（使用roleID作为前缀）
		arcfile = filepath.Join(arcPath, roleID+"-"+dialogueName+".jsonl")

		// 检查文件是否已存在
		if _, err := os.Stat(arcfile); err == nil {
			return fmt.Errorf("对话文件已存在: %s", arcfile)
		}

		// 读取角色开场白并设置到配置中
		roleDir := filepath.Join(rolePath, roleID+".role")
		prologuePath := filepath.Join(roleDir, "prologue")
		if content, err := os.ReadFile(prologuePath); err == nil && len(content) > 0 {
			c.Character.Prologue.Content = string(content)
			c.Character.Prologue.Enabled = true
		}

		// 注意：不在此处创建空文件，等到有实际对话内容时再创建
		// 初始化对话
		req = aitalkJson.NewChat(c)

	} else if arcName != newDialogueName && arcName != "" {
		// 2. 读取已有对话逻辑

		arcfile = filepath.Join(arcPath, arcName)
		// 读取存档
		req, err = aitalkJson.LoadChat(c, arcfile)
		if err != nil {
			return err
		}
		// 自动关闭开场白
		c.Character.Prologue.Enabled = false

		// 从存档文件名中提取角色ID
		fileName := filepath.Base(arcfile)
		// 解析 "roleID-dialogueName.jsonl" 格式
		roleID = extractRoleIDFromFilename(fileName)

		// 从角色目录读取显示名
		roleDir := filepath.Join(rolePath, roleID+".role")
		valuesPath := filepath.Join(roleDir, "values.json")
		if data, err := os.ReadFile(valuesPath); err == nil {
			var meta struct {
				Name string `json:"name"`
			}
			if err := json.Unmarshal(data, &meta); err == nil && meta.Name != "" {
				roleName = meta.Name
			} else {
				roleName = roleID
			}
		} else {
			roleName = roleID
		}
	}

	// 开始循环对话
	err = loopScan(req, c, arcfile, rolePath, roleID, roleName)
	if err != nil {
		return err
	}
	return nil
}

// extractRoleIDFromFilename 从对话文件名提取角色ID
// 文件名格式: roleID-dialogueName.jsonl
func extractRoleIDFromFilename(filename string) string {
	// 去掉 .jsonl 后缀
	name := filename
	if len(name) > 6 && name[len(name)-6:] == ".jsonl" {
		name = name[:len(name)-6]
	}

	// 找到第一个 "-"，前面是 roleID
	for i := 0; i < len(name); i++ {
		if name[i] == '-' {
			return name[:i]
		}
	}

	// 如果没有 "-"，返回整个名字
	return name
}
