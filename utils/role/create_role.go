package role

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"aitalk/utils/hash"
)

// RoleMeta 角色元数据
type RoleMeta struct {
	Name string `json:"name"`
}

func CreateRole(path string, name string, textEditor string) (string, error) {
	// 生成唯一ID
	roleID, err := hash.GenerateID()
	if err != nil {
		return "", fmt.Errorf("生成角色ID失败: %w", err)
	}

	// 创建角色目录（使用ID作为目录名）
	roleDir := filepath.Join(path, roleID+".role")
	if err := os.MkdirAll(roleDir, 0755); err != nil {
		return "", err
	}

	// 创建 values.json 元数据文件
	meta := RoleMeta{Name: name}
	metaPath := filepath.Join(roleDir, "values.json")
	metaData, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		os.RemoveAll(roleDir)
		return "", fmt.Errorf("序列化元数据失败: %w", err)
	}
	if err := os.WriteFile(metaPath, metaData, 0644); err != nil {
		os.RemoveAll(roleDir)
		return "", fmt.Errorf("写入元数据文件失败: %w", err)
	}

	// 创建 prologue 文件
	prologuePath := filepath.Join(roleDir, "prologue")
	if err := createFileWithEditor(prologuePath, textEditor, "角色的开场白"); err != nil {
		// 清理角色目录
		os.RemoveAll(roleDir)
		return "", err
	}

	// 创建 setting 文件
	settingPath := filepath.Join(roleDir, "setting")
	if err := createFileWithEditor(settingPath, textEditor, "角色的设定"); err != nil {
		// 清理角色目录
		os.RemoveAll(roleDir)
		return "", err
	}

	return roleID, nil
}

func createFileWithEditor(filePath string, textEditor string, description string) error {
	// 创建空文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	file.Close()

	// 如果指定了文本编辑器，使用编辑器打开文件
	if textEditor != "" {
		cmd := exec.Command(textEditor, filePath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Printf("正在使用 %s 编辑 %s...\n", textEditor, description)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("编辑器执行失败: %v", err)
		}
	} else {
		// 如果没有指定文本编辑器，提示用户
		fmt.Printf("请手动编辑文件 %s 来设置 %s\n", filePath, description)
	}

	// 检查文件内容是否为空
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}
	trimmedContent := strings.TrimSpace(string(content))
	if len(trimmedContent) == 0 {
		return fmt.Errorf("%s 不应为空", description)
	}

	return nil
}
