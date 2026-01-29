package chat

import (
	"aitalk/archive"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const newRoleName string = "创建一个新角色 / create a new role"
const newDialogueName string = "创建一个新对话 / create a new dialogue"

// 询问使用已有对话还是创建新对话
// rolePath 用于查询角色名称
func askUseWhichChatToStart(arcDir string, rolePath string) (string, error) {

	fmt.Println("** 已有的对话id清单 / dialogue id list **")

	dialogues, err := archive.ReadDialogueFiles(arcDir)
	if err != nil {
		return "", err
	}

	dialogues[0] = newDialogueName

	for k, v := range dialogues {

		if k == 0 {
			fmt.Printf("[%d] %s\n", k, v)
			continue
		}

		// 从文件名提取对话显示名
		name := strings.TrimSuffix(v, filepath.Ext(v))

		// 使用临时拷贝不要影响到name防止后续读不到hash
		curDialogueName := name
		// 清理文件名中的 hash 前缀
		if idx := strings.Index(curDialogueName, "-"); idx != -1 {
			curDialogueName = curDialogueName[idx+1:]
		}

		// 从文件名提取角色ID
		roleID := extractRoleIDFromFilename(v)

		// 查询角色名称
		roleName := getRoleNameByID(rolePath, roleID)
		if roleName != "" {
			fmt.Printf("[%d] %s (角色 / role: %s)\n", k, curDialogueName, roleName)
		} else {
			fmt.Printf("[%d] %s\n", k, curDialogueName)
		}
	}

	in := bufio.NewScanner(os.Stdin)
	fmt.Print("输入想要开始的对话id / select the number: ")
	if !in.Scan() {
		return "", fmt.Errorf("no input was given")
	}
	number, err := strconv.Atoi(in.Text())
	if err != nil {
		return "", fmt.Errorf("请输入一个已有的数字")
	}
	if chatName, ok := dialogues[number]; ok {
		return chatName, nil
	} else {
		return "", fmt.Errorf("请输入一个已有的数字")
	}
}

// getRoleNameByID 根据角色ID获取角色名称
func getRoleNameByID(rolePath string, roleID string) string {
	if roleID == "" {
		return ""
	}

	valuesPath := filepath.Join(rolePath, roleID+".role", "values.json")
	data, err := os.ReadFile(valuesPath)
	if err != nil {
		return ""
	}

	var meta struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return ""
	}

	return meta.Name
}

// 询问使用已有角色还是创建新角色
// 返回角色ID（hash值）和角色显示名
func askUseWhichRole(roleDir string) (string, string, error) {

	fmt.Println("** 已有的角色id清单 / role id list **")

	roles, err := archive.ReadRoleDirs(roleDir)
	if err != nil {
		return "", "", err
	}

	// 显示角色列表
	for k, v := range roles {
		fmt.Printf("[%d] %s (%s)\n", k, v.Name, v.ID)
	}
	fmt.Printf("[0] %s\n", newRoleName)

	in := bufio.NewScanner(os.Stdin)
	fmt.Print("输入想要开始对话的角色id / select the number: ")
	if !in.Scan() {
		return "", "", fmt.Errorf("no input was given")
	}
	number, err := strconv.Atoi(in.Text())
	if err != nil {
		return "", "", fmt.Errorf("请输入一个已有的数字")
	}

	if number == 0 {
		return "", "", nil // 创建新角色
	}

	if roleInfo, ok := roles[number]; ok {
		return roleInfo.ID, roleInfo.Name, nil
	} else {
		return "", "", fmt.Errorf("请输入一个已有的数字")
	}
}
