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
		client:    &http.Client{Timeout: 5 * time.Minute},
		messageCh: make(chan *Message),
	}
}

func (b *Bot) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer w.WriteHeader(http.StatusOK)

		var u Update
		_ = json.NewDecoder(r.Body).Decode(&u)
		b.messageCh <- &u.Message
	}
}

func (b *Bot) Messages() <-chan *Message {
	return b.messageCh
}

// SetWebhook assigns bot's webhook URL with the given URL.
func (b *Bot) SetWebhook(webhook string) error {
	params := url.Values{}
	params.Set("url", webhook)

	var r response
	err := b.sendCommand(nil, "setWebhook", params, &r)
	if err != nil {
		return err
	}

	if !r.OK {
		return fmt.Errorf("%v (%v)", r.Desc, r.ErrCode)
	}

	return nil
}

// DeleteWebhook removes webhook integration.
// TODO(ig):
func (b *Bot) deleteWebhook() error {
	return nil
}

// GetWebhookInfo retrieves current webook status.
// TODO(ig):
func (b *Bot) getWebhookInfo() (webhookinfo, error) {
	return webhookinfo{}, nil
}

// SendMessage sends text message to the recipient.
func (b *Bot) SendMessage(recipient int64, message string, opts ...SendOption) (Message, error) {
	const method = "sendMessage"
	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(recipient, 10))
	params.Set("text", message)
	mapSendOptions(&params, opts...)

	var r struct {
		response
		Message Message `json:"result"`
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
	panic("TODO")
}

// SendPhoto sends given photo to recipient. Only remote URLs are supported for now.
// A trivial example is:
//
//  b := bot.New("your-token-here")
//  photo := bot.Photo{URL: "http://i.imgur.com/6S9naG6.png"}
//  b.SendPhoto(recipient, photo, "sample image")
func (b *Bot) SendPhoto(recipient int64, photo Photo, opts ...SendOption) (Message, error) {
	const method = "sendPhoto"
	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(recipient, 10))
	params.Set("caption", photo.Caption)

	mapSendOptions(&params, opts...)
	var r struct {
		response
		Message Message `json:"result"`
	}

	var err error
	if photo.Exists() {
		params.Set("photo", photo.FileID)
		err = b.sendCommand(nil, method, params, &r)
	} else if photo.URL != "" {
		params.Set("photo", photo.URL)
		err = b.sendCommand(nil, method, params, &r)
	} else {
		err = b.sendFile(method, photo.File, "photo", params, &r)
	}

	if err != nil {
		return Message{}, err
	}

	if !r.OK {
		return Message{}, fmt.Errorf("%v (%v)", r.Desc, r.ErrCode)
	}

	return r.Message, nil
}

// SendAudio sends audio files, if you want Telegram clients to display
// them in the music player. audio must be in the .mp3 format and must not
// exceed 50 MB in size.
func (b *Bot) SendAudio(recipient int64, audio Audio, opts ...SendOption) (Message, error) {
	const method = "sendAudio"
	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(recipient, 10))
	params.Set("caption", audio.Caption)

	mapSendOptions(&params, opts...)
	var r struct {
		response
		Message Message `json:"result"`
	}

	var err error
	if audio.Exists() {
		params.Set("audio", audio.FileID)
		err = b.sendCommand(nil, method, params, &r)
	} else if audio.URL != "" {
		params.Set("audio", audio.URL)
		err = b.sendCommand(nil, method, params, &r)
	} else {
		err = b.sendFile(method, audio.File, "audio", params, &r)
	}

	if err != nil {
		return Message{}, err
	}

	if !r.OK {
		return Message{}, fmt.Errorf("%v (%v)", r.Desc, r.ErrCode)
	}

	return r.Message, nil
}

// SendDocument sends general files. Documents must not exceed 50 MB in size.
func (b *Bot) sendDocument(recipient int64, document Document, opts ...SendOption) (Message, error) {
	const method = "sendDocument"
	panic("TODO")
}

//SendSticker sends stickers with .webp extensions.
func (b *Bot) sendSticker(recipient int64, sticker Sticker, opts ...SendOption) (Message, error) {
	const method = "sendSticker"
	panic("TODO")
}

// SendVideo sends video files. Telegram clients support mp4 videos (other
// formats may be sent as Document). Video files must not exceed 50 MB in size.
func (b *Bot) sendVideo(recipient int64, video Video, opts ...SendOption) (Message, error) {
	const method = "sendVideo"
	panic("TODO")
}

// SendVoice sends audio files, if you want Telegram clients to display
// the file as a playable voice message. For this to work, your audio must be
// in an .ogg file encoded with OPUS (other formats may be sent as Audio or
// Document). audio must not exceed 50 MB in size.
func (b *Bot) sendVoice(recipient int64, audio Audio, opts ...SendOption) (Message, error) {
	const method = "sendVoice"
	panic("TODO")
}

func (b *Bot) sendVideoNote(recipient int64, videonote VideoNote, opts ...SendOption) (Message, error) {
	const method = "sendVideoNote"
	panic("TODO")
}

// SendLocation sends location point on the map.
func (b *Bot) SendLocation(recipient int64, location Location, opts ...SendOption) (Message, error) {
	const method = "sendLocation"
	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(recipient, 10))
	params.Set("latitude", strconv.FormatFloat(location.Lat, 'f', -1, 64))
	params.Set("longitude", strconv.FormatFloat(location.Long, 'f', -1, 64))

	mapSendOptions(&params, opts...)

	var r struct {
		response
		Message Message `json:"result"`
	}
	err := b.sendCommand(nil, method, params, &r)
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
	const method = "sendVenue"
	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(recipient, 10))
	params.Set("latitude", strconv.FormatFloat(venue.Location.Lat, 'f', -1, 64))
	params.Set("longitude", strconv.FormatFloat(venue.Location.Long, 'f', -1, 64))
	params.Set("title", venue.Title)
	params.Set("address", venue.Address)

	mapSendOptions(&params, opts...)

	var r struct {
		response
		Message Message `json:"result"`
	}
	err := b.sendCommand(nil, method, params, &r)
	if err != nil {
		return Message{}, err
	}

	if !r.OK {
		return Message{}, fmt.Errorf("%v (%v)", r.Desc, r.ErrCode)
	}
	return r.Message, nil
}

func (b *Bot) sendContact(recipient int64, contact Contact, opts ...SendOption) (Message, error) {
	const method = "sendContact"
	panic("TODO")
}

// SendChatAction broadcasts type of action to recipient, such as `typing`,
// `uploading a photo` etc.
//
// Use this method when you need to tell the user that something is happening
// on the bot's side. The status is set for 5 seconds or less (when a message
// arrives from your bot, Telegram clients clear its typing status).
//
// Example: The ImageBot needs some time to process a request
// and upload the image. Instead of sending a text message along the lines of
// “Retrieving image, please wait…”, the bot may use SendChatAction with action
// = UploadingPhoto. The user will see a “sending photo” status for the bot.
func (b *Bot) SendChatAction(recipient int64, action ChatAction) error {
	const method = "sendChatAction"
	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(recipient, 10))
	params.Set("action", string(action))

	var r response
	err := b.sendCommand(nil, method, params, &r)
	if err != nil {
		return err

	}
	if !r.OK {
		return fmt.Errorf("%v (%v)", r.Desc, r.ErrCode)
	}

	return nil
}

// GetFile retrieves basic info about a file and prepare it for downloading.
// For the moment, bots can download files of up to 20MB in size.
// It is guaranteed that the link will be valid for at least 1 hour. When the
// link expires, a new one can be requested by calling getFile again.
func (b *Bot) GetFile(fileID string) (File, error) {
	const method = "getFile"
	params := url.Values{}
	params.Set("file_id", fileID)

	var r struct {
		response
		File File `json:"result"`
	}
	err := b.sendCommand(nil, method, params, &r)
	if err != nil {
		return File{}, err
	}

	if !r.OK {
		return File{}, fmt.Errorf("%v (%v)", r.Desc, r.ErrCode)
	}

	u := "https://api.telegram.org/file/bot" + b.token + "/" + r.File.FilePath
	r.File.URL = u

	return r.File, nil
}

// Use this method to set a new profile photo for the chat. Photos can't be
// changed for private chats. The bot must be an administrator in the chat for
// this to work and must have the appropriate admin rights.
func (b *Bot) setChatPhoto(recipient int64, photo Photo) error {
	const method = "setChatPhoto"
	panic("TODO")
}

func (b *Bot) setChatTitle(recipient int64, title string) error {
	const method = "setChatTitle"
	panic("TODO")
}

func (b *Bot) setChatDescription(recipient int64, desc string) error {
	const method = "setChatDescription"
	panic("TODO")
}

func (b *Bot) getChat(recipient int64) (Chat, error) {
	const method = "getChat"
	panic("TODO")
}

// DeleteMessage deletes a message, including service messages, with the following limitations:
// - A message can only be deleted if it was sent less than 48 hours ago.
// - Bots can delete outgoing messages in groups and supergroups.
// - Bots granted can_post_messages permissions can delete outgoing messages in channels.
// - If the bot is an administrator of a group, it can delete any message there.
// - If the bot has can_delete_messages permission in a supergroup or a channel, it can delete any message there.
func (b *Bot) deleteMessage(recipient int64, messageid int64) error {
	const method = "deleteMessage"
	panic("TODO")
}

// Use this method to delete a chat photo. Photos can't be changed for private
// chats. The bot must be an administrator in the chat for this to work and
// must have the appropriate admin rights.
func (b *Bot) deleteChatPhoto(recipient int64) error {
	const method = "setChatPhoto"
	panic("TODO")
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

// WithDisableNotification sends the message silently. Users will receive a
// notification with no sound.
func WithDisableNotification(disable bool) SendOption {
	return func(o *sendOptions) {
		o.disableNotification = disable
	}
}

func (b *Bot) getMe() (User, error) {
	var r struct {
		response
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(&v)
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(&v)
}

func mapSendOptions(m *url.Values, opts ...SendOption) {
	var o sendOptions
	for _, opt := range opts {
		if opt != nil {
			opt(&o)
		}
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

	if o.replyMarkup.Keyboard != nil {
		kb, _ := json.Marshal(o.replyMarkup)
		m.Set("reply_markup", string(kb))
	}
}

// response is a common response structure.
type response struct {
	OK      bool   `json:"ok"`
	Desc    string `json:"description"`
	ErrCode int    `json:"error_code"`
}
