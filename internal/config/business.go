package config

import (
	"fmt"
	"os"
	"strings"

	apccfg "github.com/ethereal3x/apc/config"
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

// LoadOptions 业务配置加载选项
type LoadOptions struct {
	Path string
}

// Load 从与 apc 相同的路径规则加载 business 配置段
func Load(opts LoadOptions) error {
	path := apccfg.ResolveConfigPath(apccfg.LoadOptions{Path: opts.Path})
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read business config %s: %w", path, err)
	}
	return loadFromBytes(data)
}

// loadFromBytes 从 YAML 字节解析 business 配置段
func loadFromBytes(data []byte) error {
	var file businessFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return fmt.Errorf("unmarshal business config: %w", err)
	}
	bizCfg = file.Business
	return nil
}

// Validate 校验 business 配置必填项与格式
func Validate(cfg BusinessConfig) error {
	if strings.TrimSpace(cfg.JWT.Secret) == "" {
		return fmt.Errorf("business.jwt.secret is required")
	}
	keyLen := len(cfg.SecretKey.Encryption)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return fmt.Errorf("business.secret_key.encryption must be 16, 24 or 32 bytes, got %d", keyLen)
	}
	return nil
}

// GetBusinessConfig 获取 business 配置
func GetBusinessConfig() BusinessConfig {
	return bizCfg
}
