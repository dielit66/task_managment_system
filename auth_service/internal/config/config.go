package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	LogLevel int8 `yaml:"log_level"`
	Server   struct {
		JWTSecret string `yaml:"jwt_secret"`
		Port      string `yaml:"port"`
	} `yaml:"server"`
	Databse struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		Host     string `yaml:"host"`
	} `yaml:"database"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yaml", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
