package json

import (
	"aitalk/config"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func NewChat(c *config.Config) *ChatReq {
	// 初始化消息列表，系统设定必加
	messages := []Message{
		{Role: "system", Content: c.Character.CharacterSetting},
	}

	// 如果启用了开场白，添加 assistant 的开场消息
	if c.Character.Prologue.Enabled {
		messages = append(messages, Message{
			Role:    "assistant",
			Content: c.Character.Prologue.Content,
		})
	}

	return &ChatReq{
		Model:     c.ModelApi.Model,
		Messages:  messages,
		Thinking:  Thinking{Type: c.ModelApi.Thinking},
		Stream:    c.ModelApi.Stream,
		MaxTokens: c.ModelApi.MaxTokens,
		Temp:      c.ModelApi.Temp,
	}
}

func LoadChat(c *config.Config, arcFilePath string) (*ChatReq, error) {
	arc, err := loadMessagesFromFile(arcFilePath)
	if err != nil {
		return nil, err
	}
	if arc == nil {
		return nil, fmt.Errorf("no data in %s", arcFilePath)
	}
	return &ChatReq{
		Model:     c.ModelApi.Model,
		Messages:  arc,
		Thinking:  Thinking{Type: c.ModelApi.Thinking},
		Stream:    c.ModelApi.Stream,
		MaxTokens: c.ModelApi.MaxTokens,
		Temp:      c.ModelApi.Temp,
	}, nil
}

// 读取 .jsonl 文件并返回 []Message
func loadMessagesFromFile(path string) ([]Message, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var msgs []Message
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 { // 跳过空行
			continue
		}
		var m Message
		if err := json.Unmarshal(line, &m); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, sc.Err()
}

// // 读取 JSON 文件并解析到 []Message
// func loadMessagesFromFile(path string) ([]Message, error) {
// 	data, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var msgs []Message
// 	if err := json.Unmarshal(data, &msgs); err != nil {
// 		return nil, err
// 	}
// 	return msgs, nil
// }
