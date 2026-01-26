package json

import "encoding/json"

// ToJSON 直接把 Go 值变成字符串
func TransToAiNeedJSON(req *ChatReq) (string, error) {
	b, err := json.MarshalIndent(req, "", "    ")
	return string(b), err
}
