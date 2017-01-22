package bot

import "github.com/igungor/tlbot"

type Bot struct {
	*tlbot.Bot
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

	bot := tlbot.New(cfg.Token)
	return &Bot{
		Config: cfg,
		Bot:    &bot,
		Store:  store,
	}, nil
}
