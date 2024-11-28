package model

import (
	"database/sql"
	"time"
)

type Role string

const (
	ADMIN     Role = "ADMIN"
	MODERATOR Role = "MODERATOR"
	GUEST     Role = "GUEST"
)

type Member struct {
	ID        string    `json:"id"`
	Role      Role      `json:"role"`
	UserID    string    `json:"user_id"`
	ServerID  string    `json:"server_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MemberRepository interface {
	CreateMemberTx(tx *sql.Tx, member Member) (*Member, error)
}
