package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/igungor/ilberbot"
	_ "github.com/igungor/ilberbot/commands"
)

const botname = "ilberbot"

var (
	reImage = regexp.MustCompile(`(https?:\/\/.*\.(?:png|jpg|jpeg|gif|PNG|JPG|JPEG|GIF))`)
)

// flags
var (
	debug = flag.Bool("debug", true, "enable debug")
	port  = flag.String("port", "1985", "port")
)

func printdebug(format string, args ...interface{}) {
	if *debug {
		fmt.Printf(format, args...)
	}
}

func asciifold(s string) string {
	s = strings.ToLower(s)
	r := strings.NewReplacer("ç", "c", "ğ", "g", "ı", "i", "ö", "o", "ş", "s", "ü", "u")

	return r.Replace(s)
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}()

	var u ilberbot.TelegramUpdate
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Printf("decode error: %v\n", err)
		return
	}

	printdebug("message coming from: %#v\n", u.Message.From.Username)

	if u.Message.Text == "" {
		log.Printf("no message text: %#v\n", u)
		return
	}

	keywords := strings.Fields(u.Message.Text)
	command := asciifold(keywords[0])
	chatID := u.Message.Chat.ID

	if strings.HasSuffix(command, "@"+botname) {
		command = strings.TrimSuffix(command, "@"+botname)
	}

	var args []string
	if len(keywords) > 1 {
		args = keywords[1:]
	}

	printdebug("command: %v | args: %v\n", command, args)

	result := ilberbot.Dispatch(command, args...)

	if reImage.MatchString(result) {
		go ilberbot.SetAction(chatID, "upload_photo")
		if err := ilberbot.SendPhoto(chatID, result); err != nil {
			log.Printf("Error while sending photo: %v\n", err)
			return
		}
		return
	}

	go ilberbot.SetAction(chatID, "typing")
	ilberbot.SendMessage(chatID, result)
}

func main() {
	log.SetPrefix("ilberbot: ")
	log.SetFlags(0)

	flag.Parse()

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(net.JoinHostPort("0.0.0.0", *port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
