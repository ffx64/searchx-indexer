package discord

import "time"

type Message struct {
	ID        string     `json:"id"`
	Content   string     `json:"content"`
	Timestamp time.Time  `json:"timestamp"`
	EditedAt  *time.Time `json:"edited_at"`

	// User information
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	IsBot    bool   `json:"is_bot"` 

	// Channel/Group information
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	GuildID     string `json:"guild_id"`
	GuildName   string `json:"guild_name"`

	// Additional metadata
	IsPinned    bool         `json:"is_pinned"`
	IsTTS       bool         `json:"is_tts"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
}

type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Color       int    `json:"color"`
}
