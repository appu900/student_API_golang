package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Address string `yaml:"address" env:"ADDRESS" env-required:"true" env-default:"localhost:5000"`
}

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true"`
	HttpServer  `yaml:"http_server"`
}

func MustLoad() *Config {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		flags := flag.String("config", "", "path to configuration file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
			log.Fatal("config path is required")
		}
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file not found")
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("config: %v", err.Error())
	}

	return &cfg
}
