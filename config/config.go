package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ModelApi   ModelApi  `yaml:"modelApi"`   // 对应 YAML 节点 modelApi
	Character  Character `yaml:"character"`  // 角色卡配置
	TextEditor string    `yaml:"textEditor"` // 文本编辑器
	Player     Player    `yaml:"player"`     // 玩家本人设定
}

type Character struct {
	CharacterSetting string   `yaml:"characterSetting"`
	Prologue         Prologue `yaml:"prologue"`
	Memory           int      `yaml:"memory"`
}

type Player struct {
	Name string `yaml:"name"`
}

type Prologue struct {
	Enabled bool   `yaml:"enabled"`
	Content string `yaml:"content"`
}

type ModelApi struct {
	Url       string  `yaml:"url"`       // 请求地址
	Model     string  `yaml:"model"`     // 对应 YAML 节点 model
	ApiKey    string  `yaml:"apiKey"`    // 对应 YAML 节点 apiKey
	Thinking  string  `yaml:"thinking"`  // 是否启用思考模式
	Stream    bool    `yaml:"stream"`    // 是否启用流式输出
	MaxTokens int     `yaml:"maxTokens"` // 最大生成token数
	Temp      float32 `yaml:"temp"`      // 温度参数（随机性）
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
