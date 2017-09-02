package bot

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Token        string `json:"token"`
	Webhook      string `json:"webhook"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	Debug        bool   `json:"debug"`
	DatabasePath string `json:"dbPath"`

	GoogleAPIKey         string `json:"googleAPIKey"`
	GoogleSearchEngineID string `json:"googleSearchEngineID"`
	OpenweathermapAppID  string `json:"openWeatherMapAppID"`
}

func readConfig(configPath string) (*Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("token field can not be empty")
	}
	if cfg.Webhook == "" {
		return nil, fmt.Errorf("webhook field can not be empty")
	}
	return &cfg, nil
}
