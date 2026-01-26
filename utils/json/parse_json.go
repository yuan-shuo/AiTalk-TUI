package json

import "encoding/json"

// --------------- 解析函数 ---------------
func ParseResponse(resp string, data *ChatResponse) error {
	err := json.Unmarshal([]byte(resp), data)
	return err
}
