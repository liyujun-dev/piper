package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Profiles       []Profile `yaml:"profiles"`
	CurrentProfile string    `yaml:"current-profile"`
}

type Profile struct {
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

func AddProfile(cfg *Config, profile Profile) error {
	// 检查配置文件是否已经存在
	for _, p := range cfg.Profiles {
		if p.Name == profile.Name {
			return fmt.Errorf("profile '%s' already exists", profile.Name)
		}
	}
	cfg.Profiles = append(cfg.Profiles, profile)
	return nil
}

func RemoveProfile(cfg *Config, profileName string) error {
	for i, p := range cfg.Profiles {
		if p.Name == profileName {
			cfg.Profiles = append(cfg.Profiles[:i], cfg.Profiles[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("profile '%s' not found", profileName)
}

func ListProfiles(cfg *Config) []Profile {
	return cfg.Profiles
}
