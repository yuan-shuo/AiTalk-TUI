package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// 全局变量用于存储文件路径
var (
	ConfigPath  string
	ArchivePath string
	RolePath    string
)

var rootCmd = &cobra.Command{
	Use:   "cmda",
	Short: "AI Talk Command Line Tool",
	Long:  "A command line tool for AI talk",
}

func Execute(configPath, archivePath, rolePath string) error {
	// 设置文件路径
	ConfigPath = configPath
	ArchivePath = archivePath
	RolePath = rolePath

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(log.Writer(), err)
		log.Fatal(err)
	}
	return nil
}
