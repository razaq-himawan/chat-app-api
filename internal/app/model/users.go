package model

import (
	"time"

	"github.com/coder/websocket"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	ImageURL  string    `json:"image_url"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	CreateUser(user User) (*User, error)
	FindUserByField(field, value string) (*User, error)
}

type UserService interface {
	RegisterUser(registerPayload UserRegisterPayload) (*User, error)
	LoginUser(u *User) (string, error)
	CheckUserCredentials(loginPayload UserLoginPayload) (*User, error)
	CheckIfEmailOrUsernameExists(registerPayload UserRegisterPayload) error
	GetUserByID(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByUsername(username string) (*User, error)
}

type UserRegisterPayload struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=130"`
	Name     string `json:"name" validate:"required"`
	ImageURL string `json:"image_url,omitempty" validate:"omitempty,url"`
	Email    string `json:"email" validate:"required,email"`
}

type UserLoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

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
