package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/igungor/ilber/command"
	"github.com/igungor/tlbot"
)

// flags
var (
	flagConfig = flag.String("c", "./ilber.conf", "configuration file path")
)

func usage() {
	fmt.Fprintf(os.Stderr, "ilber is a multi-purpose Telegram bot\n\n")
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "  ilber -c <path of ilber.conf>\n\n")
	fmt.Fprintf(os.Stderr, "flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetPrefix("ilber: ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Usage = usage
	flag.Parse()

	config, err := readConfig(*flagConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "configuration error: %v\n", err)
		os.Exit(1)
	}

	b := tlbot.New(config.Token)
	if err := b.SetWebhook(config.Webhook); err != nil {
		log.Fatalf("error while setting webhook: %v", err)
	}
	log.Printf("Webhook set to %v\n", config.Webhook)

	if config.Profile {
		go func() {
			log.Println("Exposing profile information on http://:6969")
			log.Printf("profile error: %v", http.ListenAndServe(":6969", nil))
		}()
	}

	ctx := newCtxWithValues(config)

	messages := b.Listen(net.JoinHostPort(config.Host, config.Port))
	for msg := range messages {
		log.Printf("%v\n", msg)

		// react only to user sent messages
		if msg.IsService() {
			continue
		}
		// is message a bot command?
		cmdname := msg.Command()
		if cmdname == "" {
			continue
		}

		// is the command even registered?
		cmd := command.Lookup(cmdname)
		if cmd == nil {
			continue
		}

		// it is. cool, run it!
		go cmd.Run(ctx, &b, &msg)
	}
}

type config struct {
	Token   string `json:"token"`
	Webhook string `json:"webhook"`
	Host    string `json:"host"`
	Port    string `json:"port"`
	Debug   bool   `json:"debug"`
	Profile bool   `json:"profile"`

	GoogleAPIKey         string `json:"googleAPIKey"`
	GoogleSearchEngineID string `json:"googleSearchEngineID"`
	OpenweathermapAppID  string `json:"openWeatherMapAppID"`
}

func readConfig(configpath string) (config *config, err error) {
	f, err := os.Open(configpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}
	if config.Token == "" {
		return nil, fmt.Errorf("token field can not be empty")
	}
	if config.Webhook == "" {
		return nil, fmt.Errorf("webhook field can not be empty")
	}
	return config, nil
}

func newCtxWithValues(c *config) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "googleAPIKey", c.GoogleAPIKey)
	ctx = context.WithValue(ctx, "googleSearchEngineID", c.GoogleSearchEngineID)
	ctx = context.WithValue(ctx, "openWeatherMapAppID", c.OpenweathermapAppID)
	return ctx
}
