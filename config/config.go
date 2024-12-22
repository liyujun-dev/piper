package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Contexts       []Context `yaml:"contexts"`
	CurrentContext string    `yaml:"current-context"`
}

type Context struct {
	Name     string `yaml:"name"`
	Provider string `yaml:"provider"`
	Token    string `yaml:"token"`
	Server   string `yaml:"server"`
}

func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(filePath string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func AddContext(cfg *Config, context Context) error {
	// 检查上下文是否已经存在
	for _, ctx := range cfg.Contexts {
		if ctx.Name == context.Name {
			return fmt.Errorf("context '%s' already exists", context.Name)
		}
	}
	cfg.Contexts = append(cfg.Contexts, context)
	return nil
}

func RemoveContext(cfg *Config, contextName string) error {
	for i, ctx := range cfg.Contexts {
		if ctx.Name == contextName {
			cfg.Contexts = append(cfg.Contexts[:i], cfg.Contexts[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("context '%s' not found", contextName)
}

func ListContexts(cfg *Config) []Context {
	return cfg.Contexts
}
