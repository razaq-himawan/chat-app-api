package model

import (
	"database/sql"
	"time"
)

type ChannelType string

const (
	TEXT  ChannelType = "TEXT"
	AUDIO ChannelType = "AUDIO"
	VIDEO ChannelType = "VIDEO"
)

type Channel struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Type      ChannelType `json:"channel_type"`
	UserID    string      `json:"user_id"`
	ServerID  string      `json:"server_id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type ChannelRepository interface {
	CreateChannelTx(tx *sql.Tx, channel Channel) (*Channel, error)
}
