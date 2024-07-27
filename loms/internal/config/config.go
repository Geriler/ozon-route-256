package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ApplicationName string         `yaml:"application_name"`
	Env             string         `yaml:"env"`
	GRPC            AddressConfig  `yaml:"grpc"`
	HTTP            AddressConfig  `yaml:"http"`
	Database        DatabaseConfig `yaml:"database"`
	TimeoutStop     time.Duration  `yaml:"timeout_stop"`
	Tracer          AddressConfig  `yaml:"tracer"`
	Kafka           KafkaConfig    `yaml:"kafka"`
	Outbox          OutboxConfig   `yaml:"outbox"`
}

type AddressConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DatabaseConfig struct {
	DSNs []string `yaml:"dsn"`
}

type KafkaConfig struct {
	Addresses               []string      `yaml:"addresses"`
	Topic                   string        `yaml:"topic"`
	ProducerMessageInterval time.Duration `yaml:"producer_message_interval"`
}

type OutboxConfig struct {
	ClearTableInterval time.Duration `yaml:"clear_table_interval"`
	OldDataDuration    time.Duration `yaml:"old_data_duration"`
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
