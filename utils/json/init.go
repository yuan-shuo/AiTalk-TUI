package json

import "aitalk/config"

func NewReqStruct(c *config.Config) *ChatReq {
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
