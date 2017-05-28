package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/ilber/command"
)

// flags
var (
	flagConfig = flag.String("c", "./ilber.conf", "configuration file path")
)

func usage() {
	fmt.Fprintf(os.Stderr, "ilber is a multi-purpose Telegram bot\n\n")
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "  ilber -c path-to-ilber.conf\n\n")
	fmt.Fprintf(os.Stderr, "flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetPrefix("ilber: ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Usage = usage
	flag.Parse()

	b, err := bot.New(*flagConfig)
	if err != nil {
		log.Fatalf("error initializing the bot: %v\n", err)
	}

	err = b.SetWebhook(b.Config.Webhook)
	if err != nil {
		log.Fatalf("error while setting webhook: %v", err)
	}
	log.Printf("Webhook set to %v\n", b.Config.Webhook)

	http.HandleFunc("/", b.Handler())

	go func() {
		addr := net.JoinHostPort(b.Config.Host, b.Config.Port)
		log.Fatal(http.ListenAndServe(addr, nil))
	}()

	if b.Config.Profile {
		go func() {
			log.Println("Exposing profile information on http://0.0.0.0:6969")
			log.Printf("profile error: %v", http.ListenAndServe(":6969", nil))
		}()
	}

	ctx := context.Background()
	for msg := range b.Messages() {
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
		go cmd.Run(ctx, b, msg)
	}
}
