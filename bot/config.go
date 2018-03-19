package bot

import (
	"fmt"
	"net"

	"github.com/burntsushi/toml"
)

type Config struct {
	Token        string `toml:"token"`
	Webhook      string `toml:"webhook"`
	Addr         string `toml:"addr"`
	Debug        bool   `toml:"debug"`
	DatabasePath string `toml:"database-path"`

	GoogleAPIKey         string `toml:"google-api-key"`
	GoogleSearchEngineID string `toml:"google-search-engine-id"`
	OpenweathermapAppID  string `toml:"openweathermap-app-id"`
	AlphaVantageToken    string `toml:"alphavantage-token"`
}

func (c Config) validate() error {
	if c.Addr == "" {
		return fmt.Errorf("addr must be specified (host:port)")
	}
	_, _, err := net.SplitHostPort(c.Addr)
	if err != nil {
		return fmt.Errorf("invalid 'addr' specified")
	}

	if c.Token == "" {
		return fmt.Errorf("'token' must be specified")
	}

	if c.Webhook == "" {
		return fmt.Errorf("'webhook' must be specified")
	}

	return nil
}

func Load(path string) (Config, error) {
	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("could not decode config file: %v", err)
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
