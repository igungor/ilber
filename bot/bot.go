package bot

import (
	"fmt"
	"log"

	"github.com/igungor/telegram"
	"github.com/kelseyhightower/envconfig"
)

type Bot struct {
	*telegram.Bot
	Config Config
	Logger *log.Logger
}

func New(logger *log.Logger) (*Bot, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	bot := telegram.New(cfg.Token)
	return &Bot{
		Config: cfg,
		Bot:    bot,
		Logger: logger,
	}, nil
}

type Config struct {
	Token string `required:"true"`
	Debug bool   `default:"false"`

	GoogleAPIKey         string `required:"false" envconfig:"google_api_key"`
	GoogleSearchEngineID string `required:"false" envconfig:"google_search_engine_id"`
	OpenweathermapAppID  string `required:"false" envconfig:"openweathermap_app_id"`
}

func Load() (Config, error) {
	var cfg Config
	err := envconfig.Process("ilber", &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("process config: %v", err)
	}
	return cfg, nil
}
