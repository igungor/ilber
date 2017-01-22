package telegram

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"
)

type Action string

// Types of actions to broadcast
const (
	Typing            Action = "typing"
	UploadingPhoto    Action = "upload_photo"
	UploadingVideo    Action = "upload_video"
	UploadingAudio    Action = "upload_audio"
	UploadingDocument Action = "upload_document"
	FindingLocation   Action = "find_location"
)

// User represents a Telegram user or bot.
type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Chat represents a Telegram chat.
type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// IsGroupChat reports whether the message is originally sent from a chat group.
func (c Chat) IsGroupChat() bool { return c.Type == "group" }

type Update struct {
	ID      int64   `json:"update_id"`
	Payload Message `json:"message"`
}

// Message represents a message to be sent.
type Message struct {
	// Unique message identifier
	ID int64 `json:"message_id"`

	// Sender (optional. can be empty for messages sent to channel)
	From User `json:"from"`

	// Date is when the message was sent in Unix time
	Unixtime int64 `json:"date"`

	// Conversation the message belongs to — user in case of a private chat,
	// group in case of a group chat
	Chat Chat `json:"chat"`

	// For forwarded messages, sender of the original message (Optional)
	ForwardFrom User `json:"forward_from"`

	// For forwarded messages, date the original message was sent in
	// Unix time (Optional)
	ForwardDate int64 `json:"forward_date"`

	// For replies, the original message. Note that the Message
	// object in this field will not contain further reply_to_message fields
	// even if it itself is a reply (Optional)
	ReplyTo *Message `json:"reply_to_message"`

	// For text messages, the actual UTF-8 text of the message (Optional)
	Text string `json:"text"`

	// Message is an audio file, information about the file (Optional)
	Audio Audio `json:"audio"`

	// Message is a general file, information about the file (Optional)
	Document Document `json:"document"`

	// Message is a photo, available sizes of the photo (Optional)
	Photos []Photo `json:"photo"`

	// Message is a sticker, information about the sticker (Optional)
	Sticker Sticker `json:"sticker"`

	// Message is a video, information about the video (Optional)
	Video Video `json:"video"`

	// Message is a shared contact, information about the contact (Optional)
	Contact Contact `json:"contact"`

	// Message is a shared location, information about the location (Optional)
	Location Location `json:"location"`

	// A new member was added to the group, information about them
	// (this member may be bot itself) (Optional)
	JoinedUser User `json:"new_chat_participant"`

	// A member was removed from the group, information about them
	// (this member may be bot itself) (Optional)
	LeftUser User `json:"left_chat_participant"`

	// A group title was changed to this value (Optional)
	NewChatTitle string `json:"new_chat_title"`

	// A group photo was change to this value (Optional)
	NewChatPhoto []Photo `json:"new_chat_photo"`

	// Informs that the group photo was deleted (Optional)
	ChatPhotoDeleted bool `json:"delete_chat_photo"`

	// Informs that the group has been created (Optional)
	GroupChatCreated bool `json:"group_chat_created"`
}

// String returns a human-readable representation of Message.
func (m Message) String() string {
	var buf bytes.Buffer
	if m.Chat.IsGroupChat() {
		buf.WriteString(fmt.Sprintf(`From user("%v %v [%v - %v] in Group(%v [%v])" `, m.From.FirstName, m.From.LastName, m.From.Username, m.From.ID, m.Chat.Title, m.Chat.ID))
	} else {
		buf.WriteString(fmt.Sprintf(`From user("%v %v [%v - %v])" `, m.From.FirstName, m.From.LastName, m.From.Username, m.From.ID))
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

type File struct {
	// File is embedded in most of the types. So a `File` prefix is used
	FileID   string `json:"file_id"`
	FileSize int    `json:"file_size"`
	FilePath string `json:"file_path"`

	Name string    `json:"-"`
	Body io.Reader `json:"-"`
	URL  string    `json:"-"`
}

// Exists reports whether the file is already at Telegram servers.
func (f File) Exists() bool { return f.FileID != "" }

// Photo represents one size of a photo or a file/sticker thumbnail.
type Photo struct {
	File
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Caption string `json:"caption"`
}

// Audio represents an audio file.
type Audio struct {
	File
	Duration  int    `json:"duration"`
	Performer string `json:"performer"`
	Title     string `json:"title"`
	MimeType  string `json:"mime_type"`
}

// Document represents a general file (as opposed to photos and audio files).
type Document struct {
	File
	Filename  string `json:"file_name"`
	Thumbnail Photo  `json:"thumb"`
	MimeType  string `json:"mime_type"`
}

// Sticker represents a sticker.
type Sticker struct {
	File
	Width     int   `json:"width"`
	Height    int   `json:"height"`
	Thumbnail Photo `json:"thumb"`
}

// Video represents a video file.
type Video struct {
	File
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Duration  int    `json:"duration"`
	Thumbnail Photo  `json:"thumb"`
	MimeType  string `json:"mime_type"`
	Caption   string `json:"caption"`
}

// Voice represents an voice note.
type Voice struct {
	File
	Duration int    `json:"duration"`
	MimeType string `json:"mime_type"`
}

// Location represents a point on the map.
type Location struct {
	Lat  float64 `json:"latitude"`
	Long float64 `json:"longitude"`
}

// Venue represents a venue
type Venue struct {
	Location     Location `json:"location"`
	Title        string   `json:"title"`
	Address      string   `json:"address"`
	FoursquareID string   `json:"foursquare_id"`
}

// Contact represents a phone contact.
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserID      string `json:"user_id"`
}

type ReplyMarkup struct {
	Keyboard   [][]string `json:"keyboard,omitempty"`
	Resize     bool       `json:"resize_keyboard,omitempty"`
	OneTime    bool       `json:"one_time_keyboard,omitempty"`
	Selective  bool       `json:"selective,omitempty"`
	Hide       bool       `json:"hide_keyboard,omitempty"`
	ForceReply bool       `json:"force_reply,omitempty"`
}
