package cmd

import (
	"aitalk/archive"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const newChatName string = "START A NEW DIALOGUE"

// 询问使用新对话还是读取指定存档
func askUseWhichChatToStart(arcDir string) (string, error) {

	fmt.Println("** input the number to start that dialogue **")

	hash, err := archive.ReadHistoryFiles(arcDir)
	if err != nil {
		return "", err
	}

	hash[0] = newChatName

	for k, v := range hash {
		fmt.Printf("[%d] %s\n", k, v)
	}

	in := bufio.NewScanner(os.Stdin)
	fmt.Print("select the number: ")
	if !in.Scan() {
		return "", fmt.Errorf("no input was given")
	}
	number, err := strconv.Atoi(in.Text())
	if chatName, ok := hash[number]; ok {
		return chatName, nil
	} else {
		return "", fmt.Errorf("input number is wrong, please check again")
	}
}
