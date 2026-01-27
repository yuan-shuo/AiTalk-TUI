package cmd

import (
	"aitalk/config"
	"aitalk/utils/json"
	"path/filepath"
)

func Run(c *config.Config, arcPath string) error {

	var req *json.ChatReq

	// 询问读取什么存档
	arcName, err := askUseWhichChatToStart(arcPath)
	if err != nil {
		return err
	}
	if arcName == newChatName {
		// 开启新对话
		req = json.NewChat(c)
	} else if arcName != newChatName && arcName != "" {
		// 读取存档
		req, err = json.LoadChat(c, filepath.Join(arcPath, arcName))
		if err != nil {
			return err
		}
		// 自动关闭开场白
		c.Character.Prologue.Enabled = false
	}

	// 开始循环对话
	err = loopScan(req, c)
	if err != nil {
		return err
	}
	return nil
}
