package json

import (
	"encoding/json"
	"os"
)

// 追加一条消息
func AppendMessage(path string, m Message) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	return enc.Encode(m) // 自带换行
}
