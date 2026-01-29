package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"aitalk/archive"
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
		roles, err := archive.ReadRoleDirs(RolePath)
		if err != nil {
			fmt.Printf("读取角色目录失败：%v\n", err)
			return
		}

		// 显示角色列表 [编号] 名称 (ID)
		for k, v := range roles {
			fmt.Printf("[%d] %s (%s)\n", k, v.Name, v.ID)
		}

		// 选择角色
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("请输入角色编号： ")
		var roleIndex int
		if scanner.Scan() {
			roleIndex, _ = strconv.Atoi(scanner.Text())
		}

		roleInfo, ok := roles[roleIndex]
		if !ok {
			fmt.Println("角色编号无效！")
			return
		}

		// 选择要编辑的内容
		fmt.Printf("\n角色：%s (%s)\n", roleInfo.Name, roleInfo.ID)
		fmt.Println("请选择要编辑的内容：")
		fmt.Println("[0] 角色名称 (name)")
		fmt.Println("[1] 开场白 (prologue)")
		fmt.Println("[2] 人物设定 (character setting)")
		fmt.Print("请输入选择： ")
		var editChoice int
		if scanner.Scan() {
			editChoice, _ = strconv.Atoi(scanner.Text())
		}

		// 确定要编辑的文件路径
		roleDir := filepath.Join(RolePath, roleInfo.ID+".role")
		
		switch editChoice {
		case 0:
			// 编辑角色名称
			fmt.Printf("当前角色名称：%s\n", roleInfo.Name)
			fmt.Print("请输入新的角色名称： ")
			var newName string
			if scanner.Scan() {
				newName = scanner.Text()
			}
			if newName == "" {
				fmt.Println("角色名称不能为空！")
				return
			}
			
			// 更新 values.json
			valuesPath := filepath.Join(roleDir, "values.json")
			meta := struct {
				Name string `json:"name"`
			}{Name: newName}
			
			metaData, err := json.MarshalIndent(meta, "", "  ")
			if err != nil {
				fmt.Printf("序列化元数据失败：%v\n", err)
				return
			}
			
			if err := os.WriteFile(valuesPath, metaData, 0644); err != nil {
				fmt.Printf("保存角色名称失败：%v\n", err)
				return
			}
			fmt.Println("角色名称修改成功！")
			
		case 1:
			// 编辑开场白
			filePath := filepath.Join(roleDir, "prologue")
			editFile(filePath, c.TextEditor)
			
		case 2:
			// 编辑人物设定
			filePath := filepath.Join(roleDir, "setting")
			editFile(filePath, c.TextEditor)
			
		default:
			fmt.Println("选择无效！")
			return
		}
	},
}

// editFile 使用指定编辑器编辑文件
func editFile(filePath string, textEditor string) {
	// 确保文件存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("创建文件失败：%v\n", err)
			return
		}
	}

	// 使用配置文件中的文本编辑器打开文件
	editorCmd := exec.Command(textEditor, filePath)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr
	if err := editorCmd.Run(); err != nil {
		fmt.Printf("打开编辑器失败：%v\n", err)
		return
	}

	fmt.Println("编辑完成！")
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
