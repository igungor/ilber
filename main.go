package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/ilber/command"
)

func usage() {
	fmt.Fprintf(os.Stderr, "ilber is a multi-purpose Telegram bot\n\n")
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "  ilber -c path-to-ilber.conf\n\n")
	fmt.Fprintf(os.Stderr, "flags:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	logger := log.New(os.Stdout, "ilber: ", log.LstdFlags|log.Lshortfile)

	b, err := bot.NewFromEnv(logger)
	if err != nil {
		logger.Fatalf("Could not initialize the bot: %v\n", err)
	}

	http.HandleFunc("/", b.Handler())

	go func() {
		b.Logger.Fatal(http.ListenAndServe(b.Config.Addr, nil))
	}()

	ctx := context.Background()
	for msg := range b.Messages() {
		b.Logger.Printf("%v\n", msg)

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

		go cmd.Run(ctx, b, msg)
	}
}
