package json

import "aitalk/config"

func NewReqStruct(c *config.Config) *ChatReq {
	// 请求体
	return &ChatReq{
		Model: c.ModelApi.Model,
		Messages: []Message{
			{Role: "system", Content: "你是一名AI助理，可以为用户提供帮助"},
			// {Role: "assistant", Content: "hi!"},
			// {Role: "user", Content: "who are you?"},
		},
		Thinking:  Thinking{Type: "disabled"},
		Stream:    false,
		MaxTokens: 65536,
		Temp:      1.0,
	}
}
