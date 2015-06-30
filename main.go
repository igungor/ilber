package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	token   = os.Getenv("ILBERBOT_TOKEN")
	baseURL = "https://api.telegram.org/bot" + token
)

// flags
var (
	debug = flag.Bool("debug", true, "enable debug")
)

type command func(args ...string) string

var commandMap = map[string]command{}

func register(name string, command command) {
	if _, ok := commandMap[name]; ok {
		log.Println("panic: command '%s' is already registered", name)
		return
	}

	commandMap[name] = command
}

func dispatch(command string, args ...string) string {
	cmd, ok := commandMap[command]
	if !ok {
		log.Printf("command '%s' not found", command)
		return ""
	}

	return cmd(args...)
}

// Telegram Bot Response
type (
	Update struct {
		UpdateID int `json:"update_id"`
		Message  Message
	}

	Message struct {
		From User
		Chat GroupChat
		Date int
		Text string
	}

	User struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
	}

	GroupChat struct {
		ID    int `json:"id"`
		Title string
	}
)

func sendMessage(chatID int, text string) {
	u, _ := url.Parse(baseURL + "/sendMessage")
	v := u.Query()
	v.Set("chat_id", strconv.Itoa(chatID))
	v.Set("text", text)
	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	printdebug("sendMessage status: %v\n", resp.Status)
}

func setAction(chatID int, action string) {
	u, _ := url.Parse(baseURL + "/sendChatAction")
	v := u.Query()
	v.Set("chat_id", strconv.Itoa(chatID))
	v.Set("action", action)
	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	printdebug("setAction status: %v\n", resp.Status)
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}()

	var u Update
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Printf("decode error: %v\n", err)
		return
	}

	printdebug("incoming message: %#v\n", u)

	if u.Message.Text == "" {
		log.Printf("no incoming message text: %#v\n", u)
		return
	}

	keywords := strings.Fields(u.Message.Text)
	command := asciifold(keywords[0])

	var args []string
	if len(keywords) > 1 {
		args = keywords[1:]
	}

	printdebug("command: %v | args: %v\n", command, args)

	result := dispatch(command, args...)

	// show typing message whiling fetching movies
	setAction(u.Message.Chat.ID, "typing")

	sendMessage(u.Message.Chat.ID, result)
}

func main() {
	log.SetPrefix("ilberbot: ")
	log.SetFlags(0)

	flag.Parse()

	if token == "" {
		log.Fatal("ILBERBOT_TOKEN must be set")
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":1985", nil)
}
