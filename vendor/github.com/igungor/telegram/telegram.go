package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ParseMode string

// Parse modes
const (
	ModeNone     ParseMode = ""
	ModeMarkdown ParseMode = "Markdown"
	ModeHTML     ParseMode = "HTML"
)

// Bot represent a Telegram bot.
type Bot struct {
	token     string
	baseURL   string
	client    *http.Client
	messageCh chan *Message
}

// New creates a new Telegram bot with the given token, which is given by
// Botfather. See https://core.telegram.org/bots#botfather
func New(token string) *Bot {
	return &Bot{
		token:     token,
		baseURL:   fmt.Sprintf("https://api.telegram.org/bot%v/", token),
		client:    &http.Client{Timeout: 30 * time.Second},
		messageCh: make(chan *Message),
	}
}

func (b *Bot) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer w.WriteHeader(http.StatusOK)

		var u Update
		_ = json.NewDecoder(r.Body).Decode(&u)
		b.messageCh <- &u.Payload
	}
}

func (b *Bot) Messages() <-chan *Message {
	return b.messageCh
}

// SetWebhook assigns bot's webhook url with the given url.
func (b *Bot) SetWebhook(webhook string) error {
	params := url.Values{}
	params.Set("url", webhook)

	var r struct {
		OK      bool   `json:"ok"`
		Desc    string `json:"description"`
		ErrCode int    `json:"errorcode"`
	}
	err := b.sendCommand(nil, "setWebhook", params, &r)
	if err != nil {
		return err
	}

	if !r.OK {
		return fmt.Errorf("%v (%v)", r.Desc, r.ErrCode)
	}

	return nil
}

// SendMessage sends text message to the recipient. Callers can send plain
// text or markdown messages by setting mode parameter.
func (b *Bot) SendMessage(recipient int64, message string, opts ...SendOption) (Message, error) {
	params := url.Values{
		"chat_id": {strconv.FormatInt(recipient, 10)},
		"text":    {message},
	}

	mapSendOptions(&params, opts...)

	var r struct {
		OK      bool   `json:"ok"`
		Desc    string `json:"description"`
		ErrCode int    `json:"error_code"`
		Message Message
	}
	err := b.sendCommand(nil, "sendMessage", params, &r)
	if err != nil {
		return r.Message, err
	}

	if !r.OK {
		return Message{}, fmt.Errorf("%v (%v)", r.Desc, r.ErrCode)
	}
	return r.Message, nil
}

func (b *Bot) forwardMessage(recipient User, message Message) (Message, error) {
	panic("not implemented yet")
}

// SendPhoto sends given photo to recipient. Only remote URLs are supported for now.
// A trivial example is:
//
//  b := bot.New("your-token-here")
//  photo := bot.Photo{URL: "http://i.imgur.com/6S9naG6.png"}
//  err := b.SendPhoto(recipient, photo, "sample image", nil)
func (b *Bot) SendPhoto(recipient int64, photo Photo, opts ...SendOption) (Message, error) {
	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(recipient, 10))
	params.Set("caption", photo.Caption)

	mapSendOptions(&params, opts...)
	var r struct {
		OK      bool    `json:"ok"`
		Desc    string  `json:"description"`
		ErrCode int     `json:"error_code"`
		Message Message `json:"message"`
	}

	var err error
	if photo.Exists() {
		params.Set("photo", photo.FileID)
		err = b.sendCommand(nil, "sendPhoto", params, &r)
	} else if photo.URL != "" {
		params.Set("photo", photo.URL)
		err = b.sendCommand(nil, "sendPhoto", params, &r)
	} else {
		err = b.sendFile("sendPhoto", photo.File, "photo", params, &r)
	}

	if err != nil {
		return Message{}, err
	}

	if !r.OK {
		return Message{}, fmt.Errorf("%v (%v)", r.Desc, r.ErrCode)
	}

	return r.Message, nil
}

func (b *Bot) sendFile(method string, f File, form string, params url.Values, v interface{}) error {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile(form, f.Name)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, f.Body)
	if err != nil {
		return err
	}

	for k, v := range params {
		w.WriteField(k, v[0])
	}

	err = w.Close()
	if err != nil {
		return err
	}

	resp, err := b.client.Post(b.baseURL+method, w.FormDataContentType(), &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(&v)
}

// SendAudio sends audio files, if you want Telegram clients to display
// them in the music player. audio must be in the .mp3 format and must not
// exceed 50 MB in size.
func (b *Bot) sendAudio(recipient int64, audio Audio, opts ...SendOption) (Message, error) {
	panic("not implemented yet")
}

// SendDocument sends general files. Documents must not exceed 50 MB in size.
func (b *Bot) sendDocument(recipient int64, document Document, opts ...SendOption) (Message, error) {
	panic("not implemented yet")
}

//SendSticker sends stickers with .webp extensions.
func (b *Bot) sendSticker(recipient int64, sticker Sticker, opts ...SendOption) (Message, error) {
	panic("not implemented yet")
}

// SendVideo sends video files. Telegram clients support mp4 videos (other
// formats may be sent as Document). Video files must not exceed 50 MB in size.
func (b *Bot) sendVideo(recipient int64, video Video, opts ...SendOption) (Message, error) {
	panic("not implemented yet")
}

// SendVoice sends audio files, if you want Telegram clients to display
// the file as a playable voice message. For this to work, your audio must be
// in an .ogg file encoded with OPUS (other formats may be sent as Audio or
// Document). audio must not exceed 50 MB in size.
func (b *Bot) sendVoice(recipient int64, audio Audio, opts ...SendOption) (Message, error) {
	panic("not implemented yet")
}

// SendLocation sends location point on the map.
func (b *Bot) SendLocation(recipient int64, location Location, opts ...SendOption) (Message, error) {
	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(recipient, 10))
	params.Set("latitude", strconv.FormatFloat(location.Lat, 'f', -1, 64))
	params.Set("longitude", strconv.FormatFloat(location.Long, 'f', -1, 64))

	mapSendOptions(&params, opts...)

	var r struct {
		OK      bool    `json:"ok"`
		Desc    string  `json:"description"`
		ErrCode int     `json:"errorcode"`
		Message Message `json:"message"`
	}
	err := b.sendCommand(nil, "sendLocation", params, &r)
	if err != nil {
		return Message{}, err
	}

	if !r.OK {
		return Message{}, fmt.Errorf("%v (%v)", r.Desc, r.ErrCode)
	}

	return r.Message, nil
}

// SendVenue sends information about a venue.
func (b *Bot) SendVenue(recipient int64, venue Venue, opts ...SendOption) (Message, error) {
	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(recipient, 10))
	params.Set("latitude", strconv.FormatFloat(venue.Location.Lat, 'f', -1, 64))
	params.Set("longitude", strconv.FormatFloat(venue.Location.Long, 'f', -1, 64))
	params.Set("title", venue.Title)
	params.Set("address", venue.Address)

	mapSendOptions(&params, opts...)

	var r struct {
		OK      bool    `json:"ok"`
		Desc    string  `json:"description"`
		ErrCode int     `json:"errorcode"`
		Message Message `json:"message"`
	}
	err := b.sendCommand(nil, "sendVenue", params, &r)
	if err != nil {
		return Message{}, err
	}

	if !r.OK {
		return Message{}, fmt.Errorf("%v (%v)", r.Desc, r.ErrCode)
	}
	return r.Message, nil
}

// SendChatAction broadcasts type of action to recipient, such as `typing`,
// `uploading a photo` etc.
func (b *Bot) SendChatAction(recipient int64, action Action) error {
	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(recipient, 10))
	params.Set("action", string(action))

	var r struct {
		OK      bool   `json:"ok"`
		Desc    string `json:"description"`
		ErrCode int    `json:"error_code"`
	}

	err := b.sendCommand(nil, "sendChatAction", params, &r)
	if err != nil {
		return err

	}
	if !r.OK {
		return fmt.Errorf("%v (%v)", r.Desc, r.ErrCode)
	}

	return nil
}

// sendOptions configure a SendMessage call. sendOptions are set by the
// SendOption values passed to SendMessage.
type sendOptions struct {
	replyTo int64

	parseMode ParseMode

	disableWebPagePreview bool

	disableNotification bool

	replyMarkup ReplyMarkup
}

// SendOption configures how we configure the message to be sent.
type SendOption func(*sendOptions)

// WithParseMode returns a SendOption which sets the message format, such as
// HTML, Markdown etc.
func WithParseMode(mode ParseMode) SendOption {
	return func(o *sendOptions) {
		o.parseMode = mode
	}
}

// WithReplyTo returns a SendOption which sets the message to be replied to.
func WithReplyTo(to int64) SendOption {
	return func(o *sendOptions) {
		o.replyTo = to
	}
}

// WithReplyMarkup returns a SendOption which configures a custom keyboard for
// the sent message.
func WithReplyMarkup(markup ReplyMarkup) SendOption {
	return func(o *sendOptions) {
		o.replyMarkup = markup
	}
}

// WithDisableWebPagePreview returns a SendOption which disables webpage
// previews if the message contains a link.
func WithDisableWebPagePreview(disable bool) SendOption {
	return func(o *sendOptions) {
		o.disableWebPagePreview = disable
	}
}

func WithDisableNotification(disable bool) SendOption {
	return func(o *sendOptions) {
		o.disableNotification = disable
	}
}

func (b *Bot) GetFile(fileID string) (File, error) {
	params := url.Values{}
	params.Set("file_id", fileID)

	var r struct {
		OK      bool   `json:"ok"`
		Desc    string `json:"description"`
		ErrCode int    `json:"errorcode"`
		File    File   `json:"result"`
	}
	err := b.sendCommand(nil, "getFile", params, &r)
	if err != nil {
		return File{}, err
	}

	if !r.OK {
		return File{}, fmt.Errorf("%v (%v)", r.Desc, r.ErrCode)
	}

	return r.File, nil
}

func (b *Bot) GetFileDownloadURL(fileID string) (string, error) {
	f, err := b.GetFile(fileID)
	if err != nil {
		return "", err
	}

	u := "https://api.telegram.org/file/bot" + b.token + "/" + f.FilePath
	return u, nil
}

func (b *Bot) sendCommand(ctx context.Context, method string, params url.Values, v interface{}) error {
	req, err := http.NewRequest("POST", b.baseURL+method, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	if ctx == nil {
		ctx = context.Background()
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(&v)
}

func (b *Bot) getMe() (User, error) {
	var r struct {
		OK      bool   `json:"ok"`
		Desc    string `json:"description"`
		ErrCode int    `json:"error_code"`

		User User `json:"result"`
	}
	err := b.sendCommand(nil, "getMe", url.Values{}, &r)
	if err != nil {
		return User{}, err
	}

	if !r.OK {
		return User{}, fmt.Errorf("%v (%v)", r.Desc, r.ErrCode)
	}

	return r.User, nil
}

func mapSendOptions(m *url.Values, opts ...SendOption) {
	var o sendOptions
	for _, opt := range opts {
		opt(&o)
	}

	if o.replyTo != 0 {
		m.Set("reply_to_message_id", strconv.FormatInt(o.replyTo, 10))
	}

	if o.disableWebPagePreview {
		m.Set("disable_web_page_preview", "true")
	}

	if o.disableNotification {
		m.Set("disable_notification", "true")
	}

	if o.parseMode != ModeNone {
		m.Set("parse_mode", string(o.parseMode))
	}

	// TODO: map ReplyMarkup options as well
}
