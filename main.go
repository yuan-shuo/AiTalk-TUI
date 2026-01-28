package main

import (
	"aitalk/cmd"
	"path/filepath"
)

func main() {
	// 定义文件路径
	configPath := filepath.Join(".", "data", "etc", "config.yaml")
	archivePath := filepath.Join(".", "data", "archive")
	rolePath := filepath.Join(".", "data", "role")

	// 传入文件路径并执行命令
	cmd.Execute(configPath, archivePath, rolePath)
}
