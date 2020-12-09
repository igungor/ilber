package bot

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

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
