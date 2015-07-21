package ilberbot

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	TelegramResult struct {
		OK          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	}

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

func SendMessage(chatID int, text string) error {
	urlvalues := url.Values{
		"chat_id": {strconv.Itoa(chatID)},
		"text":    {text},
	}

	resp, err := http.PostForm(baseURL+"/sendMessage", urlvalues)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r TelegramResult
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	if r.OK {
		return nil
	}

	return fmt.Errorf("%v (errcode: %v)", r.Description, r.ErrorCode)
}

func SetAction(chatID int, action string) error {
	// action: typing
	// action: upload_{audio,video,document}
	// action: find_location

	urlvalues := url.Values{
		"chat_id": {strconv.Itoa(chatID)},
		"action":  {action},
	}

	resp, err := http.PostForm(baseURL+"/sendChatAction", urlvalues)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r TelegramResult
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	if r.OK {
		return nil
	}

	return fmt.Errorf("%v (errcode: %v)", r.Description, r.ErrorCode)
}

func SendPhoto(chatID int, url string) error {
	// set status to "sending photo..."
	go SetAction(chatID, "upload_photo")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Fetching '%v' failed", url)
	}

	var buf bytes.Buffer

	w := multipart.NewWriter(&buf)

	part, err := w.CreateFormFile("photo", "image.jpg")
	if err != nil {
		return err
	}

	io.Copy(part, resp.Body)

	w.WriteField("chat_id", strconv.Itoa(chatID))

	contenttype := w.FormDataContentType()
	if err := w.Close(); err != nil {
		return err
	}

	resp, err = http.Post(baseURL+"/sendPhoto", contenttype, &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r TelegramResult
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	if !r.OK {
		return fmt.Errorf("%v (errcode: %v)", r.Description, r.ErrorCode)
	}

	return nil
}
