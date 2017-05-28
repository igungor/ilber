package bot

import "github.com/igungor/telegram"

type Bot struct {
	*telegram.Bot
	Config *Config
	Store  *Store
}

func New(configPath string) (*Bot, error) {
	cfg, err := readConfig(configPath)
	if err != nil {
		return nil, err
	}

	store := NewStore(cfg.DatabasePath)
	err = store.Open()
	if err != nil {
		return nil, err
	}

	bot := telegram.New(cfg.Token)
	return &Bot{
		Config: cfg,
		Bot:    bot,
		Store:  store,
	}, nil
}
