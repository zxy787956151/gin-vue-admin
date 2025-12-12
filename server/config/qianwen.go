package config

type Qianwen struct {
	ApiKey string `mapstructure:"api-key" json:"api-key" yaml:"api-key"` // 通义千问 API Key
	Model  string `mapstructure:"model" json:"model" yaml:"model"`       // 使用的模型，默认 qwen-turbo
}


