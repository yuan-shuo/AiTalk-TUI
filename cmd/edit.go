package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"aitalk/config"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit configuration or role settings",
	Long:  "Edit configuration file or role settings",
}

var editRoleCmd = &cobra.Command{
	Use:   "role",
	Short: "Edit role settings",
	Long:  "Edit role prologue or character setting",
	Run: func(cobraCmd *cobra.Command, args []string) {
		// 加载配置文件以获取文本编辑器
		c, err := config.LoadFrom(ConfigPath)
		if err != nil {
			fmt.Printf("加载配置文件失败：%v\n", err)
			return
		}

		// 列出已有角色
		fmt.Println("请选择要编辑的角色：")
		roleDirs, err := os.ReadDir(RolePath)
		if err != nil {
			fmt.Printf("读取角色目录失败：%v\n", err)
			return
		}

		roleMap := make(map[int]string)
		for i, dir := range roleDirs {
			if dir.IsDir() {
				roleName := strings.TrimSuffix(dir.Name(), ".role")
				roleMap[i+1] = roleName
				fmt.Printf("[%d] %s\n", i+1, roleName)
			}
		}

		// 选择角色
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("请输入角色编号： ")
		var roleIndex int
		if scanner.Scan() {
			roleIndex, _ = strconv.Atoi(scanner.Text())
		}

		roleName, ok := roleMap[roleIndex]
		if !ok {
			fmt.Println("角色编号无效！")
			return
		}

		// 选择要编辑的内容
		fmt.Println("请选择要编辑的内容：")
		fmt.Println("[1] 开场白 (prologue)")
		fmt.Println("[2] 人物设定 (character setting)")
		fmt.Print("请输入选择： ")
		var editChoice int
		if scanner.Scan() {
			editChoice, _ = strconv.Atoi(scanner.Text())
		}

		// 确定要编辑的文件路径
		var filePath string
		roleDir := filepath.Join(RolePath, roleName+".role")
		switch editChoice {
		case 1:
			filePath = filepath.Join(roleDir, "prologue")
		case 2:
			filePath = filepath.Join(roleDir, "setting")
		default:
			fmt.Println("选择无效！")
			return
		}

		// 确保角色目录存在
		if err := os.MkdirAll(roleDir, 0755); err != nil {
			fmt.Printf("创建角色目录失败：%v\n", err)
			return
		}

		// 确保文件存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			_, err := os.Create(filePath)
			if err != nil {
				fmt.Printf("创建文件失败：%v\n", err)
				return
			}
		}

		// 使用配置文件中的文本编辑器打开文件
		editorCmd := exec.Command(c.TextEditor, filePath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr
		if err := editorCmd.Run(); err != nil {
			fmt.Printf("打开编辑器失败：%v\n", err)
			return
		}

		fmt.Println("编辑完成！")
	},
}

var editConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Edit configuration file",
	Long:  "Edit config.yaml file",
	Run: func(cobraCmd *cobra.Command, args []string) {
		// 检查配置文件是否存在
		if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
			fmt.Println("配置文件不存在，请先运行 init 命令初始化！")
			return
		}

		// 加载配置文件以获取文本编辑器
		c, err := config.LoadFrom(ConfigPath)
		if err != nil {
			fmt.Printf("加载配置文件失败：%v\n", err)
			return
		}

		// 使用配置文件中的文本编辑器打开文件
		editorCmd := exec.Command(c.TextEditor, ConfigPath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr
		if err := editorCmd.Run(); err != nil {
			fmt.Printf("打开编辑器失败：%v\n", err)
			return
		}

		fmt.Println("编辑完成！")
	},
}

func init() {
	editCmd.AddCommand(editRoleCmd)
	editCmd.AddCommand(editConfigCmd)
	rootCmd.AddCommand(editCmd)
}
