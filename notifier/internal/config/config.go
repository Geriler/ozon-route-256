package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ApplicationName string      `yaml:"application_name"`
	Env             string      `yaml:"env"`
	Kafka           KafkaConfig `yaml:"kafka"`
}

type KafkaConfig struct {
	Addresses       []string `yaml:"addresses"`
	Topic           string   `yaml:"topic"`
	ConsumerGroupID string   `yaml:"consumer_group_id"`
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
