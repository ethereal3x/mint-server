package app

import (
	"os"
	"strings"

	apccfg "github.com/ethereal3x/apc/config"
)

const (
	envOTELAuthUsername = "OTEL_AUTH_USERNAME"
	envOTELAuthPassword = "OTEL_AUTH_PASSWORD"
)

// applyTracingEnvOverrides 从环境变量覆盖 OTLP tracing 认证配置
func applyTracingEnvOverrides(cfg *apccfg.Config) {
	if cfg == nil {
		return
	}
	reporter := &cfg.Plugin.Tracing.Reporter
	if username := strings.TrimSpace(os.Getenv(envOTELAuthUsername)); username != "" {
		reporter.Auth.Username = username
	}
	if password := os.Getenv(envOTELAuthPassword); password != "" {
		reporter.Auth.Password = password
	}
}
