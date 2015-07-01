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
		From  User
		Chat  GroupChat
		Date  int
		Text  string
		Photo []PhotoSize
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

	PhotoSize struct {
		FileID   string `json:"file_id"`
		Width    int    `json:"width"`
		Height   int    `json:"height"`
		FileSize int    `json:"file_size"`
	}
)

func sendMessage(chatID int, text string) {
	urlvalues := url.Values{
		"chat_id": {strconv.Itoa(chatID)},
		"text":    {text},
	}

	resp, err := http.PostForm(baseURL+"/sendMessage", urlvalues)
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

	urlvalues := url.Values{
		"chat_id": {strconv.Itoa(chatID)},
		"action":  {action},
	}

	resp, err := http.PostForm(baseURL+"/sendChatAction", urlvalues)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	printdebug("setAction status: %v\n", resp.Status)
}
