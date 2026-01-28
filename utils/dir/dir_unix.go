//go:build darwin || linux || freebsd

package dir

import (
	"fmt"
	"os"
	"path/filepath"
)

// getDataDir 获取Unix-like系统的应用数据目录
// macOS: ~/Library/Application Support
// Linux/FreeBSD: ~/.local/share 或 $XDG_DATA_HOME
func getDataDir() (string, error) {
	// 首先检查 XDG_DATA_HOME 环境变量（Linux标准）
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome != "" {
		return dataHome, nil
	}

	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("无法获取用户主目录: %w", err)
	}

	// 检查是否是 macOS
	if os.Getenv("GOOS") == "darwin" || filepath.Join(homeDir, "Library") != "" {
		// 检查 Library 目录是否存在来判断是否是 macOS
		libraryPath := filepath.Join(homeDir, "Library")
		if _, err := os.Stat(libraryPath); err == nil {
			// macOS: ~/Library/Application Support
			return filepath.Join(homeDir, "Library", "Application Support"), nil
		}
	}

	// Linux/FreeBSD: ~/.local/share
	return filepath.Join(homeDir, ".local", "share"), nil
}
