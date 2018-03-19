package bot

import (
	"fmt"
	"log"

	"github.com/igungor/telegram"
)

type Bot struct {
	*telegram.Bot
	Config Config
	Logger *log.Logger
}

func New(configPath string, logger *log.Logger) (*Bot, error) {
	cfg, err := Load(configPath)
	if err != nil {
		return nil, err
	}

	bot := telegram.New(cfg.Token)
	if err = bot.SetWebhook(cfg.Webhook); err != nil {
		return nil, fmt.Errorf("could not set webhook: %v", err)
	}
	return &Bot{
		Config: cfg,
		Bot:    bot,
		Logger: logger,
	}, nil
}
