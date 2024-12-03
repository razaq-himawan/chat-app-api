package model

import "github.com/coder/websocket"

type WSType string

const (
	DM      WSType = "DM"
	CHANNEL WSType = "CHANNEL"
)

type WebSocketUser struct {
	UserID         string          `json:"user_id"`
	Conn           *websocket.Conn `json:"-"`
	Type           WSType          `json:"type"`
	ConversationID string          `json:"conversation_id,omitempty"`
	ChannelID      string          `json:"channel_id,omitempty"`
	IsOnline       bool            `json:"is_online"`
}
