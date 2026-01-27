package cmd

import (
	"aitalk/config"
	"aitalk/utils/json"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func Run(c *config.Config, arcPath string) error {
	// 确保存档目录存在
	if err := os.MkdirAll(arcPath, 0755); err != nil {
		return err
	}

	var req *json.ChatReq
	var arcfile string

	// 询问读取什么存档
	arcName, err := askUseWhichChatToStart(arcPath)
	if err != nil {
		return err
	}
	if arcName == newChatName {
		// 开启新对话
		req = json.NewChat(c)

		// 询问对话名称
		in := bufio.NewScanner(os.Stdin)
		fmt.Print("enter the name of this new dialogue: ")
		if !in.Scan() {
			return fmt.Errorf("no input was given")
		}

		// 拼接存档文件路径
		arcfile = filepath.Join(arcPath, in.Text()+".jsonl")
		_, err := os.Create(arcfile)
		if err != nil {
			return err
		}

		// 写入初始化内容
		for _, v := range req.Messages {
			json.AppendMessage(arcfile, v)
		}

	} else if arcName != newChatName && arcName != "" {
		arcfile = filepath.Join(arcPath, arcName)
		// 读取存档
		req, err = json.LoadChat(c, arcfile)
		if err != nil {
			return err
		}
		// 自动关闭开场白
		c.Character.Prologue.Enabled = false
	}

	// 开始循环对话
	err = loopScan(req, c, arcfile)
	if err != nil {
		return err
	}
	return nil
}
