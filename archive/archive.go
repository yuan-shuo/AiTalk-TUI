package archive

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ReadRoleDirs 扫描 root 下所有 *.role 目录
// 返回 map[int]string，key 从 1 开始递增，value 就是目录名（如 "dog.role"）
func ReadRoleDirs(root string) (map[int]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, e := range entries {
		if e.IsDir() && strings.HasSuffix(e.Name(), ".role") {
			dirs = append(dirs, e.Name())
		}
	}

	sort.Strings(dirs)

	out := make(map[int]string, len(dirs))
	for idx, d := range dirs {
		out[idx+1] = d
	}
	return out, nil
}

// ReadDialogueFiles 扫描 dir 下所有 .jsonl 文件
// 返回 map[int]string，key 从 1 开始递增，value 为文件名（仅文件名，不含路径）。
func ReadDialogueFiles(dir string) (map[int]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// 收集所有 .jsonl 文件名
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if filepath.Ext(name) == ".jsonl" {
			files = append(files, name)
		}
	}

	// 按字典序排列，保证顺序一致
	sort.Strings(files)

	// 构建 map[int]string
	out := make(map[int]string, len(files))
	for idx, name := range files {
		out[idx+1] = name
	}
	return out, nil
}
