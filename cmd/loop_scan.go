package cmd

import (
	"aitalk/config"
	"aitalk/core"
	"aitalk/utils/json"
	"bufio"
	"fmt"
	"os"
)

func loopScan(req *json.ChatReq, c *config.Config) error {
	in := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(`YOU (enter "Ctrl+C" to exit) >>> `)
		if !in.Scan() {
			break
		}

		text := in.Text()

		aiResp, err := core.Chat(text, req, c)
		if err != nil {
			return err
		}
		fmt.Println("agent:", aiResp)
	}
	return nil
}
