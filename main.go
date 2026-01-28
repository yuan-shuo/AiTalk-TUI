package main

import (
	"aitalk/cmd"
	"aitalk/config"
	"log"
	"path/filepath"
)

func main() {
	configPath := filepath.Join(".", "data", "etc", "config.yaml")
	archivePath := filepath.Join(".", "data", "archive")
	rolePath := filepath.Join(".", "data", "role")

	// 加载配置文件
	c, err := config.LoadFrom(configPath)
	if err != nil {
		log.Fatalf("config read failed: %s", err)
	}

	// reqJson, err := json.TransToAiNeedJSON(json.NewReqStruct(c))
	// fmt.Println(reqJson)

	err = cmd.Run(c, archivePath, rolePath)
	if err != nil {
		log.Fatalf("error!: %s", err)
	}
}
