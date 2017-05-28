package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/igungor/telegram"
)

// flags
var (
	token   = flag.String("token", "", "telegram bot token")
	webhook = flag.String("webhook", "", "webhook url")
	host    = flag.String("host", "127.0.0.1", "host to listen to")
	port    = flag.String("port", "1986", "port to listen to")
)

func usage() {
	fmt.Fprintf(os.Stderr, "echobot is an echo server for testing Telegram bots\n\n")
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "  echobot -token <insert-your-telegrambot-token> -webhook <insert-your-webhook-url>\n\n")
	fmt.Fprintf(os.Stderr, "flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("echobot: ")
	flag.Usage = usage
	flag.Parse()

	if *webhook == "" {
		log.Printf("missing webhook parameter\n\n")
		flag.Usage()
	}
	if *token == "" {
		log.Printf("missing token parameter\n\n")
		flag.Usage()
	}

	bot := telegram.New(*token)
	err := bot.SetWebhook(*webhook)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", bot.Handler())

	go func() {
		log.Fatal(http.ListenAndServe(net.JoinHostPort(*host, *port), nil))
	}()

	for msg := range bot.Messages() {
		go func(msg *telegram.Message) {
			// echo the message as *bold*
			txt := "*" + msg.Text + "*"
			_, err := bot.SendMessage(msg.Chat.ID, txt, nil)
			if err != nil {
				log.Printf("Error while sending message. Err: %v\n", err)
			}
		}(msg)
	}
}
