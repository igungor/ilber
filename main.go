package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/igungor/ilberbot/command"
	"github.com/igungor/tlbot"
)

// flags
var (
	token   = flag.String("token", "", "telegram bot token")
	webhook = flag.String("webhook", "", "webhook url")
	host    = flag.String("host", "", "host to listen to")
	port    = flag.String("port", "1985", "port to listen to")
)

func usage() {
	fmt.Fprintf(os.Stderr, "ilberbot is a multi-purpose Telegram bot\n\n")
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "  ilberbot -token <insert-your-telegrambot-token> -webhook <insert-your-webhook-url>\n\n")
	fmt.Fprintf(os.Stderr, "flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("ilberbot: ")
	flag.Usage = usage
	flag.Parse()

	if *token == "" {
		log.Printf("missing token parameter\n\n")
		flag.Usage()
	}
	if *webhook == "" {
		log.Printf("missing webhook parameter\n\n")
		flag.Usage()
	}

	b := tlbot.New(*token)
	err := b.SetWebhook(*webhook)
	if err != nil {
		log.Fatal(err)
	}

	messages := b.Listen(net.JoinHostPort(*host, *port))
	for msg := range messages {
		// is message a command?
		cmdname := msg.Command()
		if cmdname == "" {
			continue
		}

		// does it even exist?
		cmd := command.Lookup(cmdname)
		if cmd == nil {
			continue
		}

		// cool, run it!
		cmd.Run(&b, &msg)
	}
}
