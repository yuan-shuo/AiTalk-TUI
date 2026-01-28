package main

import (
	"aitalk/cmd"
	"aitalk/config"
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cmda",
	Short: "AI Talk Command Line Tool",
	Long:  "A command line tool for AI talk",
}

var talkCmd = &cobra.Command{
	Use:   "talk",
	Short: "Start an AI talk session",
	Long:  "Start an AI talk session with a character",
	Run: func(cobraCmd *cobra.Command, args []string) {
		configPath := filepath.Join(".", "data", "etc", "config.yaml")
		archivePath := filepath.Join(".", "data", "archive")
		rolePath := filepath.Join(".", "data", "role")

		// 加载配置文件
		c, err := config.LoadFrom(configPath)
		if err != nil {
			log.Fatalf("config read failed: %s", err)
		}

		err = cmd.Run(c, archivePath, rolePath)
		if err != nil {
			log.Fatalf("error!: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(talkCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(log.Writer(), err)
		log.Fatal(err)
	}
}
