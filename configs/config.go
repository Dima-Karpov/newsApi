package configs

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	RSS           []string `json:"rss"`
	RequestPeriod int      `json:"request_period"`
}

// LoadConfig загружает конфигурацию из файла
func LoadConfig(filePath string) (Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, errors.New("config file not found")
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, errors.New("config file parse error")
	}

	return cfg, nil
}
