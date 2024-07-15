package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ApplicationName string        `yaml:"application_name"`
	Env             string        `yaml:"env"`
	HTTP            AddressConfig `yaml:"http"`
	Product         ProductConfig `yaml:"product"`
	GRPC            AddressConfig `yaml:"grpc"`
	TimeoutStop     time.Duration `yaml:"timeout_stop"`
	Tracer          AddressConfig `yaml:"tracer"`
}

type ProductConfig struct {
	BaseUrl  string `yaml:"base_url"`
	Token    string `yaml:"token"`
	RPSLimit int    `yaml:"rps_limit"`
}

type AddressConfig struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout" omitempty:"true"`
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

	err := validateConfig(&cfg)
	if err != nil {
		panic("invalid config: " + err.Error())
	}

	return cfg
}

func validateConfig(c *Config) error {
	if c.Product.RPSLimit <= 0 {
		return fmt.Errorf("RPSLimit must be greater than 0")
	}

	return nil
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
