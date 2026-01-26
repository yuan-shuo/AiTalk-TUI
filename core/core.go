package core

import (
	"aitalk/client"
	"aitalk/config"
	"aitalk/utils/json"
	"fmt"
)

func Chat(text string, req *json.ChatReq, c *config.Config) (string, error) {
	// 请求体转换为json
	req.Messages = append(req.Messages, json.Message{Role: "user", Content: text})
	reqJson, err := json.TransToAiNeedJSON(req)
	if err != nil {
		return "", fmt.Errorf("trans go data to json failed: %s", err)
	}

	// 携带参数对模型接口发起请求
	resp, err := client.PostModelApi(c.ModelApi.Url, reqJson, c.ModelApi.ApiKey)
	if err != nil {
		return "", fmt.Errorf("post failed: %s", err)
	}

	// 解析模型返回结果
	var data json.ChatResponse
	err = json.ParseResponse(resp, &data)
	if err != nil {
		return "", fmt.Errorf("parase json failed: %s", err)
	}

	aiResp := data.Choices[0].Message.Content

	req.Messages = append(req.Messages, json.Message{Role: "assistant", Content: aiResp})
	// 更新对话信息
	// req.Messages = append(req.Messages,
	// 	json.Message{Role: "user", Content: text},
	// 	json.Message{Role: "assistant", Content: aiResp},
	// )

	// 返回 ai 回复内容
	return aiResp, nil
}
