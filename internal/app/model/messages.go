package model

import "time"

type Message struct {
	ID             string    `json:"id"`
	Content        string    `json:"content"`
	MemberID       string    `json:"member_id"`
	ConversationID string    `json:"conversation_id,omitempty"`
	ChannelID      string    `json:"channel_id,omitempty"`
	Deleted        bool      `json:"deleted"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
