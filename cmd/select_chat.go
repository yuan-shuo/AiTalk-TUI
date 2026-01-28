package cmd

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

	fmt.Println("** input the number to start that dialogue **")

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
	fmt.Print("select the number: ")
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
func askUseWhichRole(roleDir string) (string, error) {

	fmt.Println("** input the number to select a role **")

	hash, err := archive.ReadRoleDirs(roleDir)
	if err != nil {
		return "", err
	}

	hash[0] = newRoleName

	for k, v := range hash {
		name := strings.TrimSuffix(v, filepath.Ext(v))
		fmt.Printf("[%d] %s\n", k, name)
	}

	in := bufio.NewScanner(os.Stdin)
	fmt.Print("select the number: ")
	if !in.Scan() {
		return "", fmt.Errorf("no input was given")
	}
	number, err := strconv.Atoi(in.Text())
	if err != nil {
		return "", fmt.Errorf("请输入一个已有的数字")
	}
	if roleName, ok := hash[number]; ok {
		return roleName, nil
	} else {
		return "", fmt.Errorf("请输入一个已有的数字")
	}
}
