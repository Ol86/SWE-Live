package config

import (
	"SWE-Live/pkg/config"
)

type AppConfig struct {
	*config.Config
}

func LoadAndValidate() (*AppConfig, error) {
	rawConfig, error := config.Load()
	if error != nil {
		return nil, error
	}

	return &AppConfig{Config: rawConfig}, nil
}
