package cmd

import (
	"aitalk/config"
	"aitalk/utils/json"
	"path/filepath"
)

func Run(c *config.Config, arcPath string) error {

	var req *json.ChatReq

	arcName, err := askUseWhichChatToStart(arcPath)
	if err != nil {
		return err
	}
	if arcName == newChatName {
		req = json.NewChat(c)
	} else if arcName != newChatName && arcName != "" {
		req, err = json.LoadChat(c, filepath.Join(arcPath, arcName))
		if err != nil {
			return err
		}
		c.Character.Prologue.Enabled = false
	}

	err = loopScan(req, c)
	if err != nil {
		return err
	}
	return nil
}
