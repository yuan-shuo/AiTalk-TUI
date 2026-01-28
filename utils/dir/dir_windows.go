//go:build windows

package dir

import (
	"fmt"
	"os"
	"path/filepath"
)

// getDataDir 获取Windows系统的应用数据目录
// 使用 %APPDATA% 环境变量（通常是 C:\Users\<用户名>\AppData\Roaming）
func getDataDir() (string, error) {
	// 首先尝试获取 APPDATA 环境变量
	appData := os.Getenv("APPDATA")
	if appData == "" {
		// 如果APPDATA未设置，尝试使用LOCALAPPDATA
		appData = os.Getenv("LOCALAPPDATA")
	}
	
	if appData == "" {
		// 如果都未设置，使用用户主目录
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("无法获取用户主目录: %w", err)
		}
		appData = filepath.Join(homeDir, "AppData", "Roaming")
	}

	return appData, nil
}
