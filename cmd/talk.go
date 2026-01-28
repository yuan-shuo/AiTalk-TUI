package cmd

import (
	"aitalk/chat"
	"aitalk/config"
	"log"

	"github.com/spf13/cobra"
)

var talkCmd = &cobra.Command{
	Use:   "talk",
	Short: "Start an AI talk session",
	Long:  "Start an AI talk session with a character",
	Run: func(cobraCmd *cobra.Command, args []string) {
		// 使用全局变量获取文件路径
		configPath := ConfigPath
		archivePath := ArchivePath
		rolePath := RolePath

		// 加载配置文件
		c, err := config.LoadFrom(configPath)
		if err != nil {
			log.Fatalf("config read failed: %s", err)
		}

		err = chat.Run(c, archivePath, rolePath)
		if err != nil {
			log.Fatalf("error!: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(talkCmd)
}
