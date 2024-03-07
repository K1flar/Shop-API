package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Env      string   `yaml:"env"`
	Server   Server   `yaml:"server"`
	DataBase DataBase `yaml:"database"`
}

type Server struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Secret string `env:"SECRET" env-required:"true"`
}

type DataBase struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"poer"`
	Name     string `yaml:"name"`
	User     string `env:"USER" env-required:"true"`
	Password string `env:"PASSWORD" env-required:"true"`
}

func New(path string) (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
