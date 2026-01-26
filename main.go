package main

import (
	"aitalk/cmd"
	"aitalk/config"
	"log"
	"path/filepath"
)

func main() {

	// 加载配置文件
	c, err := config.LoadFrom(filepath.Join(".", "etc", "config.yaml"))
	if err != nil {
		log.Fatalf("config read failed: %s", err)
	}

	// reqJson, err := json.TransToAiNeedJSON(json.NewReqStruct(c))
	// fmt.Println(reqJson)

	err = cmd.Run(c)
	if err != nil {
		log.Fatalf("error!: %s", err)
	}
}
