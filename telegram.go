package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var (
	token   = os.Getenv("ILBERBOT_TOKEN")
	baseURL = "https://api.telegram.org/bot" + token
)

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
	// action: typing
	// action: upload_{audio,video,document}
	// action: find_location

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
