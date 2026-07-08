package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Settings struct {
	APIURL     string
	Token      string
	DefaultOrg string
}

var (
	ConfigPath = ""
	ConfigName = "kit"
	ConfigType = "yaml"
)

func DefaultSettings() *Settings {
	return &Settings{
		APIURL: "https://api.github.com",
	}
}

func Init() (*viper.Viper, error) {
	v := viper.New()

	if ConfigPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		ConfigPath = filepath.Join(home, ".kit")
	}

	v.SetConfigName(ConfigName)
	v.AddConfigPath(ConfigPath)
	v.AddConfigPath(".")

	v.SetEnvPrefix("KIT")
	v.AutomaticEnv()

	v.SetDefault("api_url", DefaultSettings().APIURL)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	return v, nil
}

func Save(v *viper.Viper) error {
	if err := os.MkdirAll(ConfigPath, 0o700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	return v.WriteConfigAs(filepath.Join(ConfigPath, ConfigName+"."+ConfigType))
}

func GetConfig() (*viper.Viper, error) {
	return Init()
}

func LoadSettings() (*Settings, error) {
	v, err := Init()
	if err != nil {
		return nil, err
	}

	s := DefaultSettings()
	s.APIURL = v.GetString("api_url")
	s.Token = v.GetString("token")
	s.DefaultOrg = v.GetString("default_org")
	return s, nil
}

func SaveSettings(s *Settings) error {
	v, err := Init()
	if err != nil {
		return err
	}

	v.Set("api_url", s.APIURL)
	v.Set("token", s.Token)
	v.Set("default_org", s.DefaultOrg)

	return Save(v)
}
