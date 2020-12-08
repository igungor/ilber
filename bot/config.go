package bot

import (
	"fmt"
	"net"

	"github.com/burntsushi/toml"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Token   string `toml:"token"`
	Webhook string `toml:"webhook"`
	Addr    string `toml:"addr"`
	Debug   bool   `toml:"debug"`

	GoogleAPIKey         string `toml:"google-api-key" envconfig:"google_api_key"`
	GoogleSearchEngineID string `toml:"google-search-engine-id" envconfig:"google_search_engine_id"`
	OpenweathermapAppID  string `toml:"openweathermap-app-id" envconfig:"openweathermap_app_id"`
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

func LoadFromEnv() (Config, error) {
	var cfg Config
	err := envconfig.Process("ilber", &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("process config: %v", err)
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
