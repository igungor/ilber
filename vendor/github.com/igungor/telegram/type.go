package telegram

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"
)

// ParseMode determines the markup of the text of the message.
type ParseMode string

// Parse modes
const (
	ModeNone     ParseMode = ""
	ModeMarkdown ParseMode = "Markdown"
	ModeHTML     ParseMode = "HTML"
)

// ChatAction represents bot activity.
type ChatAction string

// Types of actions to broadcast
const (
	Typing             ChatAction = "typing"
	UploadingPhoto     ChatAction = "upload_photo"
	UploadingVideo     ChatAction = "upload_video"
	UploadingAudio     ChatAction = "upload_audio"
	UploadingDocument  ChatAction = "upload_document"
	FindingLocation    ChatAction = "find_location"
	UploadingVideoNote ChatAction = "upload_video_note"
)

// User represents a Telegram user or bot.
type User struct {
	// Unique identifier for this user or bot
	ID int64 `json:"id"`

	// User‘s or bot’s first name
	Username string `json:"username"`

	// User‘s or bot’s last name
	FirstName string `json:"first_name"`

	// User‘s or bot’s username
	LastName string `json:"last_name"`

	// IETF language tag of the user's language
	LanguageCode string `json:"language_code"`
}

// Chat represents a Telegram chat.
type Chat struct {
	// Unique identifier for this chat. This number may be greater than 32 bits and
	// some programming languages may have difficulty/silent defects in
	// interpreting it. But it is smaller than 52 bits, so a signed 64 bit integer
	// or double-precision float type are safe for storing this identifier.
	ID int64 `json:"id"`

	// Type of chat, can be either “private”, “group”, “supergroup” or “channel”
	Type string `json:"type"`

	// Title, for supergroups, channels and group chats
	Title string `json:"title"`

	// Username, for private chats, supergroups and channels if available
	Username string `json:"username"`

	// First name of the other party in a private chat
	FirstName string `json:"first_name"`

	// Last name of the other party in a private chat
	LastName string `json:"last_name"`
}

// IsGroupChat reports whether the message is originally sent from a chat group.
func (c Chat) IsGroupChat() bool { return c.Type == "group" }

type Update struct {
	// The update‘s unique identifier. Update identifiers start from a certain
	// positive number and increase sequentially. This ID becomes especially handy
	// if you’re using Webhooks, since it allows you to ignore repeated updates or
	// to restore the correct update sequence, should they get out of order.
	ID int64 `json:"update_id"`

	// New incoming message of any kind — text, photo, sticker, etc.
	Message Message `json:"message,omitempty"`

	// New version of a message that is known to the bot and was edited
	EditedMessage Message `json:"edited_message,omitempty"`
}

// Message represents a message to be sent.
type Message struct {
	// Unique message identifier
	ID int64 `json:"message_id"`

	// Sender (optional. can be empty for messages sent to channel)
	From User `json:"from,omitempty"`

	// Date is when the message was sent in Unix time
	Unixtime int64 `json:"date"`

	// Conversation the message belongs to — user in case of a private chat,
	// group in case of a group chat
	Chat Chat `json:"chat"`

	// For forwarded messages, sender of the original message
	ForwardFrom User `json:"forward_from,omitempty"`

	// For messages forwarded from a channel, information about the original channel
	ForwardFromChat Chat `json:"forward_from_chat,omitempty"`

	// For forwarded messages, date the original message was sent in
	// Unix time
	ForwardDate int64 `json:"forward_date,omitempty"`

	// For replies, the original message. Note that the Message
	// object in this field will not contain further reply_to_message fields
	// even if it itself is a reply
	ReplyTo *Message `json:"reply_to_message,omitempty"`

	// For text messages, the actual UTF-8 text of the message
	Text string `json:"text,omitempty"`

	// For text messages, special entities like usernames, URLs, bot commands, etc. that appear in the text
	Entities []MessageEntity `json:"entities,omitempty"`

	// Message is an audio file, information about the file
	Audio Audio `json:"audio,omitempty"`

	// Message is a general file, information about the file
	Document Document `json:"document,omitempty"`

	// Message is a photo, available sizes of the photo
	Photos []Photo `json:"photo,omitempty"`

	// Message is a sticker, information about the sticker
	Sticker Sticker `json:"sticker,omitempty"`

	// Message is a video, information about the video
	Video Video `json:"video,omitempty"`

	// Message is a voice message, information about the file
	Voice Voice `json:"voice,omitempty"`

	// Message is a video note, information about the video message
	VideoNote VideoNote `json:"video_note,omitempty"`

	// New members that were added to the group or supergroup and information about
	// them (the bot itself may be one of these members)
	NewChatMembers []User `json:"new_chat_members,omitempty"`

	// Caption for the document, photo or video, 0-200 characters
	Caption string `json:"caption,omitempty"`

	// Message is a shared contact, information about the contact
	Contact Contact `json:"contact,omitempty"`

	// Message is a shared location, information about the location
	Location Location `json:"location,omitempty"`

	// Message is a venue, information about the venue
	Venue Venue `json:"venue,omitempty"`

	// A new member was added to the group, information about them
	// (this member may be bot itself)
	JoinedUser User `json:"new_chat_member,omitempty"`

	// A member was removed from the group, information about them
	// (this member may be bot itself)
	LeftUser User `json:"left_chat_member,omitempty"`

	// A group title was changed to this value
	NewChatTitle string `json:"new_chat_title,omitempty"`

	// A group photo was change to this value
	NewChatPhoto []Photo `json:"new_chat_photo,omitempty"`

	// Informs that the group photo was deleted
	ChatPhotoDeleted bool `json:"delete_chat_photo,omitempty"`

	// Informs that the group has been created
	GroupChatCreated bool `json:"group_chat_created,omitempty"`
}

// String returns a human-readable representation of Message.
func (m Message) String() string {
	var buf bytes.Buffer
	if m.Chat.IsGroupChat() {
		buf.WriteString(
			fmt.Sprintf(`From user("%v %v [%v - %v] in Group(%v [%v])" `,
				m.From.FirstName,
				m.From.LastName,
				m.From.Username,
				m.From.ID,
				m.Chat.Title,
				m.Chat.ID,
			))
	} else {
		buf.WriteString(
			fmt.Sprintf(`From user("%v %v [%v - %v])" `,
				m.From.FirstName,
				m.From.LastName,
				m.From.Username,
				m.From.ID,
			))
	}
	buf.WriteString(fmt.Sprintf("Message: %q", m.Text))
	return buf.String()
}

// Command returns the command's name: the first word in the message text. If
// message text starts with a `/`, function returns the command name, or else
// empty string.
func (m Message) Command() string {
	name := m.Text
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}

	j := strings.Index(name, "/")
	if j < 0 {
		return ""
	}
	return name[j+1:]
}

// Args returns all words after the first word in the message text. First word
// is meant to be the command name and can be accessed with Command method.
func (m Message) Args() []string {
	args := strings.TrimSpace(m.Text)
	i := strings.Index(args, " ")
	if i < 0 {
		return nil
	}
	return strings.Fields(args[i+1:])
}

// IsService reports whether the message is a Telegram service message, not a
// user sent message.
func (m Message) IsService() bool {
	switch {
	case m.NewChatTitle != "":
		return true
	case len(m.NewChatPhoto) > 0:
		return true
	case m.JoinedUser != User{}:
		return true
	case m.LeftUser != User{}:
		return true
	case m.GroupChatCreated:
		return true
	case m.ChatPhotoDeleted:
		return true
	}
	return false
}

// IsReply reports whether the message is a reply to another message.
func (m Message) IsReply() bool {
	return m.ReplyTo != nil
}

// Time returns the moment of message in UTC time.
func (m Message) Time() time.Time {
	return time.Unix(m.Unixtime, 0).UTC()
}

// File represents a file ready to be downloaded. The file can be downloaded
// via the GetFile method. The URL of the file is available from URL field of
// the File. It is guaranteed that the link will be valid for at least 1 hour.
// When the link expires, a new one can be requested by calling GetFile.
type File struct {
	// Unique identifier for this file
	FileID string `json:"file_id"`

	// File size, if known
	FileSize int `json:"file_size,omitempty"`

	// File path.
	FilePath string `json:"file_path,omitempty"`

	Name string    `json:"-"`
	Body io.Reader `json:"-"`
	URL  string    `json:"-"`
}

// Exists reports whether the file is already at Telegram servers.
func (f File) Exists() bool { return f.FileID != "" }

// Photo represents one size of a photo or a file/sticker thumbnail.
type Photo struct {
	File

	// Photo width
	Width int `json:"width"`

	// Photo height
	Height int `json:"height"`

	Caption string `json:"-"`
}

// Audio represents an audio file to be treated as music by Telegram clients.
type Audio struct {
	File

	// Duration of the audio in seconds as defined by sender
	Duration int `json:"duration"`

	// Performer of the audio as defined by sender or by audio tags
	Performer string `json:"performer,omitempty"`

	// Title of the audio as defined by sender or by audio tags
	Title string `json:"title,omitempty"`

	// MIME type of the file as defined by sender
	MimeType string `json:"mime_type,omitempty"`
	Caption  string `json:"-"`
}

// Document represents a general file (as opposed to photos and audio files).
type Document struct {
	File

	// Document thumbnail as defined by sender
	Thumbnail Photo `json:"thumb,omitempty"`

	// Original filename as defined by sender
	Filename string `json:"file_name,omitempty"`

	// MIME type of the file as defined by sender
	MimeType string `json:"mime_type,omitempty"`
}

// Sticker represents a sticker.
type Sticker struct {
	File

	// Sticker width
	Width int `json:"width"`

	// Sticker height
	Height int `json:"height"`

	// Sticker thumbnail in .webp or .jpg format
	Thumbnail Photo `json:"thumb,omitempty"`

	// Emoji associated with the sticker
	Emoji string `json:"emoji,omitempty"`
}

// Video represents a video file.
type Video struct {
	File

	// Video width as defined by sender
	Width int `json:"width"`

	// Video height as defined by sender
	Height int `json:"height"`

	// Duration of the video in seconds as defined by sender
	Duration int `json:"duration"`

	// Video thumbnail
	Thumbnail Photo `json:"thumb,omitempty"`

	// Mime type of a file as defined by sender
	MimeType string `json:"mime_type,omitempty"`

	Caption string `json:"-"`
}

// Voice represents an voice note.
type Voice struct {
	File

	// Duration of the audio in seconds as defined by sender
	Duration int `json:"duration"`

	// MIME type of the file as defined by sender
	MimeType string `json:"mime_type,omitempty"`
}

// VideoNote represents a video message.
type VideoNote struct {
	File

	// Video width and height as defined by sender
	Length int `json:"length"`

	// Duration of the video in seconds as defined by sender
	Duration int `json:"duration"`

	// Video thumbnail
	Thumbnail Photo `json:"thumb"`
}

// Contact represents a phone contact.
type Contact struct {
	// Contact's phone number
	PhoneNumber string `json:"phone_number"`

	// Contact's first name
	FirstName string `json:"first_name"`

	// Contact's last name
	LastName string `json:"last_name,omitempty"`

	// Contact's user identifier in Telegram
	UserID int64 `json:"user_id,omitempty"`
}

// Location represents a point on the map.
type Location struct {
	// Longitude as defined by sender
	Long float64 `json:"longitude"`

	// Latitude as defined by sender
	Lat float64 `json:"latitude"`
}

// Venue represents a venue
type Venue struct {
	// Venue location
	Location Location `json:"location"`

	// Name of the venue
	Title string `json:"title"`

	// Address of the venue
	Address string `json:"address"`

	// Foursquare identifier of the venue
	FoursquareID string `json:"foursquare_id,omitempty"`
}

// MessageEntity represents one special entity in a text message. For example,
// hashtags, usernames, URLs, etc.
type MessageEntity struct {
	// Type of the entity. Can be mention (@username), hashtag, bot_command, url,
	// email, bold (bold text), italic (italic text), code (monowidth string), pre
	// (monowidth block), text_link (for clickable text URLs), text_mention (for
	// users without usernames)
	Type string `json:"type"`

	// Offset in UTF-16 code units to the start of the entity
	// TODO(ig):
	Offset int `json:"offset"`

	// Length of the entity in UTF-16 code units
	Length int `json:"length"`

	// For “text_link” only, url that will be opened after user taps on the text
	URL string `json:"url,omitempty"`

	// For “text_mention” only, the mentioned user
	User User `json:"user,omitempty"`
}

// ReplyMarkup represents a custom keyboard with reply options.
type ReplyMarkup struct {
	// Array of button rows, each represented by an strings
	Keyboard [][]string `json:"keyboard"`

	// Optional. Requests clients to resize the keyboard vertically for optimal
	// fit (e.g., make the keyboard smaller if there are just two rows of
	// buttons).
	Resize bool `json:"resize_keyboard,omitempty"`

	// Optional. Requests clients to hide the keyboard as soon as it's been used.
	// The keyboard will still be available, but clients will automatically display
	// the usual letter-keyboard in the chat – the user can press a special button
	// in the input field to see the custom keyboard again.
	OneTime bool `json:"one_time_keyboard,omitempty"`

	// Optional. Use this parameter if you want to show the keyboard to specific
	// users only. Targets:
	// 1) users that are @mentioned in the text of the Message object
	// 2) if the bot's message is a reply, sender of the original message.
	Selective bool `json:"selective,omitempty"`
}

// KeyboardButton represents one button of the reply keyboard. For simple text
// buttons. Optional fields are mutually exclusive.
// TODO(ig):
type keyboardButton struct {
	// Text of the button. If none of the optional fields are used, it will be
	// sent to the bot as a message when the button is pressed
	Text string `json:"text"`

	// Optional. If True, the user's phone number will be sent as a contact when
	// the button is pressed. Available in private chats only
	RequestContact bool `json:"request_contact,omitempty"`

	// Optional. If True, the user's current location will be sent when the button
	// is pressed. Available in private chats only
	RequestLocation bool `json:"request_location,omitempty"`
}

// ReplyKeyboard represent the removal of already sent keyboard markup. Upon
// receiving a message with this object, Telegram clients will remove the
// current custom keyboard and display the default letter-keyboard. By default,
// custom keyboards are displayed until a new keyboard is sent by a bot. An
// exception is made for one-time keyboards that are hidden immediately after
// the user presses a button (see ReplyMarkup).
// TODO(ig):
type replyKeyboardRemove struct {
	// Requests clients to remove the custom keyboard (user will not be able to
	// summon this keyboard; if you want to hide the keyboard from sight but keep
	// it accessible, use one_time_keyboard in ReplyMarkup)
	RemoveKeybard bool `json:"remove_keybard"`

	// Optional. Use this parameter if you want to remove the keyboard for specific
	// users only. Targets:
	// 1) users that are @mentioned in the text of the Message object
	// 2) if the bot's message is a reply, sender of the original message.
	Selective bool `json:"selective,omitempty"`
}

// TODO(ig):
type webhookinfo struct {
	// Webhook URL, may be empty if webhook is not set up
	URL string `json:"url"`

	// If a custom certificate was provided for webhook certificate checks
	HasCustomCertificate bool `json:"has_custom_certificate"`

	// Number of updates awaiting delivery
	PendingUpdateCount int `json:"pending_update_count"`

	// Unix time for the most recent error that happened when trying to deliver
	// an update via webhook
	LastErrorDate int `json:"last_error_date,omitempty"`

	// Error message in human-readable format for the most recent error that
	// happened when trying to deliver an update via webhook
	LastErrorMessage string `json:"last_error_message,omitempty"`

	// Maximum allowed number of simultaneous HTTPS connections to the webhook
	// for update delivery
	MaxConnections int `json:"max_connections,omitempty"`

	// A list of update types the bot is subscribed to. Defaults to all update types
	AllowedUpdates []string `json:"allowed_updates,omitempty"`
}
