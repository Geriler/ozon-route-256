package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env"`
	HTTP        AddressConfig `yaml:"http"`
	Product     ProductConfig `yaml:"product"`
	GRPC        AddressConfig `yaml:"grpc"`
	TimeoutStop time.Duration `yaml:"timeout_stop"`
}

type ProductConfig struct {
	BaseUrl string `yaml:"base_url"`
	Token   string `yaml:"token"`
}

type AddressConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func MustLoad() Config {
	configPath := fetchConfigPath()

	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + err.Error())
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
