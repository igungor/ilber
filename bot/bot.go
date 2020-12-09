package bot

import (
	"log"

	"github.com/igungor/telegram"
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
