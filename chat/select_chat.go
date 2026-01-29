package chat

import (
	"aitalk/archive"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const newRoleName string = "创建一个新角色 / create a new role"
const newDialogueName string = "创建一个新对话 / create a new dialogue"

// 询问使用已有对话还是创建新对话
func askUseWhichChatToStart(arcDir string) (string, error) {

	fmt.Println("** 已有的对话id清单 / dialogue id list **")

	hash, err := archive.ReadDialogueFiles(arcDir)
	if err != nil {
		return "", err
	}

	hash[0] = newDialogueName

	for k, v := range hash {
		name := strings.TrimSuffix(v, filepath.Ext(v))
		fmt.Printf("[%d] %s\n", k, name)
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
	if chatName, ok := hash[number]; ok {
		return chatName, nil
	} else {
		return "", fmt.Errorf("请输入一个已有的数字")
	}
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
