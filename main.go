package main

import (
	"aitalk/cmd"
	"aitalk/utils/dir"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// 获取应用目录
	appDirs, err := dir.GetAppDirs("aitalk-tui")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting app directories: %v\n", err)
		os.Exit(1)
	}

	// 确保目录存在
	// ConfigPath 是完整文件路径，需要提取目录
	configDir := filepath.Dir(appDirs.ConfigPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(appDirs.ArchivePath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating archive directory: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(appDirs.RolePath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating role directory: %v\n", err)
		os.Exit(1)
	}

	// 传入文件路径并执行命令
	cmd.Execute(appDirs.ConfigPath, appDirs.ArchivePath, appDirs.RolePath)
}
