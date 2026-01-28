package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize AI Talk configuration",
	Long:  "Initialize AI Talk configuration by creating a config.yaml file",
	Run: func(cobraCmd *cobra.Command, args []string) {
		// 检查config.yaml文件是否存在
		if _, err := os.Stat(ConfigPath); err == nil {
			fmt.Println("配置文件已经存在，无需初始化！")
			return
		}

		// 打印各操作系统常见文本编辑器用于提示
		fmt.Println("请选择你想要使用的文本编辑器：")
		fmt.Println("Windows常见编辑器：notepad")
		fmt.Println("Linux/macOS常见编辑器：vi, vim, nano, emacs")

		// 询问用户选择的文本编辑器
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("请输入文本编辑器名称： ")
		var textEditor string
		if scanner.Scan() {
			textEditor = scanner.Text()
		}

		// 确保配置文件目录存在
		configDir := filepath.Dir(ConfigPath)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Printf("创建配置文件目录失败：%v\n", err)
			return
		}

		// 读取模板文件
		templatePath := filepath.Join(".", "embed", "skel", "config.yaml")
		tmplContent, err := os.ReadFile(templatePath)
		if err != nil {
			fmt.Printf("读取模板文件失败：%v\n", err)
			return
		}

		// 创建模板
		tmpl, err := template.New("config").Parse(string(tmplContent))
		if err != nil {
			fmt.Printf("解析模板失败：%v\n", err)
			return
		}

		// 创建配置文件
		configFile, err := os.Create(ConfigPath)
		if err != nil {
			fmt.Printf("创建配置文件失败：%v\n", err)
			return
		}
		defer configFile.Close()

		// 填充模板并写入配置文件
		data := struct {
			TextEditor string
		}{
			TextEditor: textEditor,
		}
		if err := tmpl.Execute(configFile, data); err != nil {
			fmt.Printf("写入配置文件失败：%v\n", err)
			// 清理配置文件
			os.Remove(ConfigPath)
			return
		}

		// 使用用户指定的文本编辑器打开config.yaml进行编辑
		fmt.Println("正在打开配置文件进行编辑...")
		editorCmd := exec.Command(textEditor, ConfigPath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr
		if err := editorCmd.Run(); err != nil {
			fmt.Printf("打开编辑器失败：%v\n", err)
			// 清理配置文件
			os.Remove(ConfigPath)
			fmt.Println("配置文件已清理，请重新运行 init 命令！")
			return
		}

		fmt.Println("配置文件初始化成功！")
		fmt.Printf("配置文件路径：%s\n", ConfigPath)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
