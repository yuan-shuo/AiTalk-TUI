package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ModelApi ModelApi `yaml:"modelApi"` // 对应 YAML 节点 modelApi
}

type ModelApi struct {
	Url    string `yaml:"url"`    // 请求地址
	Model  string `yaml:"model"`  // 对应 YAML 节点 model
	ApiKey string `yaml:"apiKey"` // 对应 YAML 节点 apiKey
}

// LoadFrom 从任意路径读取并解析
func LoadFrom(path string) (*Config, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config

	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
