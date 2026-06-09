package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// BusinessConfig 包含业务相关的自定义配置
type BusinessConfig struct {
	JWT       JWTConfig    `yaml:"jwt"`
	SecretKey SecretKeyCfg `yaml:"secret_key"`
}

// JWTConfig JWT 认证配置
type JWTConfig struct {
	Secret string `yaml:"secret"`
}

// SecretKeyCfg API Key 加解密密钥配置
type SecretKeyCfg struct {
	Encryption string `yaml:"encryption"`
}

type businessFile struct {
	Business BusinessConfig `yaml:"business"`
}

var bizCfg BusinessConfig

// InitBusinessConfig 从 config.yaml 加载 business 配置段
func InitBusinessConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	var file businessFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return err
	}
	bizCfg = file.Business
	return nil
}

// GetBusinessConfig 获取 business 配置
func GetBusinessConfig() BusinessConfig {
	return bizCfg
}
