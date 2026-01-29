package hash

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateID 生成一个8字符的随机唯一ID
func GenerateID() (string, error) {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
