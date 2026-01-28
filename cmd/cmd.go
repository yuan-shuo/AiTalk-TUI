package cmd

import (
	"aitalk/config"
	"aitalk/utils/json"
	"aitalk/utils/role"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	var req *json.ChatReq
	var arcfile string
	var roleBaseName string

	// 询问使用什么对话
	arcName, err := askUseWhichChatToStart(arcPath)
	if err != nil {
		return err
	}
	if arcName == newDialogueName {
		// 开启新对话

		// 询问使用什么角色
		roleName, err := askUseWhichRole(rolePath)
		if err != nil {
			return err
		}

		// 处理角色名称
		if roleName == newRoleName {
			// 创建新角色
			in := bufio.NewScanner(os.Stdin)
			fmt.Print("请输入新角色的名称: ")
			if !in.Scan() {
				return fmt.Errorf("no input was given")
			}
			newRoleName := in.Text()

			// 调用创建角色函数
			if err := role.CreateRole(rolePath, newRoleName, c.TextEditor); err != nil {
				return err
			}

			roleBaseName = newRoleName
		} else {
			roleBaseName = strings.TrimSuffix(roleName, ".role")
		}

		// 询问对话名称
		in := bufio.NewScanner(os.Stdin)
		fmt.Print("enter the name of this new dialogue: ")
		if !in.Scan() {
			return fmt.Errorf("no input was given")
		}
		dialogueName := in.Text()

		// 拼接存档文件路径
		arcfile = filepath.Join(arcPath, roleBaseName+"-"+dialogueName+".jsonl")

		// 检查文件是否已存在
		if _, err := os.Stat(arcfile); err == nil {
			return fmt.Errorf("对话文件已存在: %s", arcfile)
		}

		// 创建存档文件
		_, err = os.Create(arcfile)
		if err != nil {
			return err
		}

		// 初始化对话
		req = json.NewChat(c)

	} else if arcName != newDialogueName && arcName != "" {
		arcfile = filepath.Join(arcPath, arcName)
		// 读取存档
		req, err = json.LoadChat(c, arcfile)
		if err != nil {
			return err
		}
		// 自动关闭开场白
		c.Character.Prologue.Enabled = false

		// 从存档文件名中提取角色名称
		fileName := filepath.Base(arcfile)
		fileNameWithoutExt := strings.TrimSuffix(fileName, ".jsonl")
		parts := strings.Split(fileNameWithoutExt, "-")
		if len(parts) > 0 {
			roleBaseName = parts[0]
		}
	}

	// 开始循环对话
	err = loopScan(req, c, arcfile, rolePath, roleBaseName)
	if err != nil {
		return err
	}
	return nil
}
