package core

import (
	"aitalk/client"
	"aitalk/config"
	"aitalk/utils/json"
	"fmt"
)

func Chat(text string, req *json.ChatReq, c *config.Config, arcFilePath string, isFirstDialogue bool, prologueContent string) (string, error) {

	// 0. 记忆长度：memory 轮对话 => 最多 2*memory 条（不含 system）
	maxPairs := c.Character.Memory
	maxMsg := 2 * maxPairs // 例如 mem = 6, 6 * 2 = 12 条

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

	// 3. 队列化：保留 system(0)，只裁剪 1…end
	if len(req.Messages) > maxMsg+1 { // +1 是因为 system 不占 memory 轮
		// 删除最早的一对（user+assistant），即索引 1 和 2
		req.Messages = append(req.Messages[:1], req.Messages[3:]...)
	}

	// 当AI回复成功后，将此轮对话写入存档文件
	// 如果是第一次对话，先写入开场白
	if isFirstDialogue && prologueContent != "" {
		json.AppendMessage(arcFilePath, json.Message{Role: "assistant", Content: prologueContent})
	}
	// 写入用户回复和AI回复
	json.AppendMessage(arcFilePath, json.Message{Role: "user", Content: text})
	json.AppendMessage(arcFilePath, json.Message{Role: "assistant", Content: aiResp})

	// 返回 ai 回复内容
	return aiResp, nil
}
