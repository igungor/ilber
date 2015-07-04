package ilberbot

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var (
	token   = os.Getenv("ILBERBOT_TOKEN")
	baseURL = "https://api.telegram.org/bot" + token
)

func init() {
	if token == "" {
		log.Fatal("ILBERBOT_TOKEN must be set")
	}
}

// Telegram Bot Response
type (
	TelegramUpdate struct {
		UpdateID int `json:"update_id"`
		Message  TelegramMessage
	}

	TelegramMessage struct {
		From  TelegramUser
		Chat  TelegramGroupChat
		Date  int
		Text  string
		Photo []TelegramPhotoSize
	}

	TelegramUser struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
	}

	TelegramGroupChat struct {
		ID    int `json:"id"`
		Title string
	}

	TelegramPhotoSize struct {
		FileID   string `json:"file_id"`
		Width    int    `json:"width"`
		Height   int    `json:"height"`
		FileSize int    `json:"file_size"`
	}
)

func SendMessage(chatID int, text string) {
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
}

func SetAction(chatID int, action string) {
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
}

func SendPhoto(chatID int, url string) {
	r, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	var buf bytes.Buffer

	w := multipart.NewWriter(&buf)

	part, err := w.CreateFormFile("photo", "image.jpg")
	if err != nil {
		log.Fatal(err)
	}

	io.Copy(part, r.Body)

	w.WriteField("chat_id", strconv.Itoa(chatID))

	contenttype := w.FormDataContentType()
	if err := w.Close(); err != nil {
		log.Fatal(err)
	}

	r, err = http.Post(baseURL+"/sendPhoto", contenttype, &buf)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
}