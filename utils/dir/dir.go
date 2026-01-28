package dir

import (
	"path/filepath"
)

// AppDirs 存储应用程序目录路径
type AppDirs struct {
	ConfigPath  string // 配置文件完整路径 (e.g., .../etc/config.yaml)
	ArchivePath string // 存档文件目录
	RolePath    string // 角色文件目录
}

// GetAppDirs 获取应用程序目录
// appName 是应用程序名称，用于创建子目录
func GetAppDirs(appName string) (*AppDirs, error) {
	// 获取系统特定的应用数据目录
	dataDir, err := getDataDir()
	if err != nil {
		return nil, err
	}

	// 构建应用特定的目录路径
	appDir := filepath.Join(dataDir, appName)

	return &AppDirs{
		ConfigPath:  filepath.Join(appDir, "etc", "config.yaml"),
		ArchivePath: filepath.Join(appDir, "archive"),
		RolePath:    filepath.Join(appDir, "role"),
	}, nil
}
