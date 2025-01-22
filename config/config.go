package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

// ConfigData представляет конфигурацию приложения.
// Это структура содержит параметры, необходимые для настройки серверного приложения.
type ConfigData struct {
	// ServerAddress определяет адрес HTTP-сервера, на котором будет работать приложение.
	ServerAddress string `json:"server_address"`
}

// isParsed отслеживает, выполнена ли обработка аргументов командной строки.
var isParsed bool

// ConfigurationInterface интерфейс, в рамках проекта используется для моков юинт тестов
type ConfigurationInterface interface {
	InitConfig() (*ConfigData, error)
}

// Configuration структура конфигурации, реализующая интерфейс ConfigurationInterface
type Configuration struct{}

// InitConfig инициализирует конфигурацию приложения.
func (cs *Configuration) InitConfig() (*ConfigData, error) {
	cfg := &ConfigData{}

	configFile := os.Getenv("CONFIG")
	if configFile == "" {
		configFile = "config.json"
	}

	file, err := os.Open(configFile)

	if err != nil {
		log.Printf("Warning: Error opening config file: %v. Using default configuration.\n", err)
	} else {
		isParsed = true
		defer file.Close()

		if err = json.NewDecoder(file).Decode(cfg); err != nil {
			log.Printf("Warning: Error decoding config file: %v. Using default configuration.\n", err)
			return nil, err
		}
	}

	if !isParsed {
		flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "Адрес HTTP-сервера")
		flag.Parse()
		isParsed = true
	}

	if ServerAddress := os.Getenv("SERVER_ADDRESS"); ServerAddress != "" {
		cfg.ServerAddress = ServerAddress
	}

	return cfg, nil
}
