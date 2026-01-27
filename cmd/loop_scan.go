package cmd

import (
	"aitalk/config"
	"aitalk/core"
	"aitalk/utils/json"
	"bufio"
	"fmt"
	"os"
)

const (
	userPrompt  string = `YOU (enter "Ctrl+C" to exit):`
	agentPrompt string = `AGENT:`
)

func loopScan(req *json.ChatReq, c *config.Config, arcFilePath string) error {
	in := bufio.NewScanner(os.Stdin)
	// 如果启用了开场白，添加 assistant 的开场消息
	if c.Character.Prologue.Enabled {
		fmt.Printf("%s %s\n\n", agentPrompt, c.Character.Prologue.Content)
	}
	for {
		fmt.Printf("%s ", userPrompt)
		if !in.Scan() {
			break
		}

		text := in.Text()

		aiResp, err := core.Chat(text, req, c, arcFilePath)
		if err != nil {
			return err
		}
		fmt.Printf("\n%s %s\n\n", agentPrompt, aiResp)
	}
	return nil
}
