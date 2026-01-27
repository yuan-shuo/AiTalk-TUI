package archive

import (
	"os"
	"path/filepath"
	"sort"
)

// ReadHistoryFiles 扫描 data/chat-history 目录下所有 .json 文件，
// 返回 map[int]string，key 从 0 开始递增，value 为文件名（仅文件名，不含路径）。
func ReadHistoryFiles(dir string) (map[int]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// 收集所有 .json 文件名
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if filepath.Ext(name) == ".json" {
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
