package model

import "time"

type ServerModel struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	InviteCode string    `json:"invite_code"`
	UserID     string    `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ServerRepository interface {
	CreateServer(server ServerModel) (*ServerModel, error)
}

type ServerService interface {
	CreateNewServer(createServerPayload CreateServerPayload, userID string) (*ServerModel, error)
}

type CreateServerPayload struct {
	Name string `json:"name" validate:"required,min=3,max=30"`
}
