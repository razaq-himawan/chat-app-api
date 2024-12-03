package repository

import (
	"database/sql"
	"fmt"

	"github.com/razaq-himawan/chat-app-api/internal/app/model"
)

type ChannelRepository struct {
	db *sql.DB
}

func NewChannelRepository(db *sql.DB) *ChannelRepository {
	return &ChannelRepository{db: db}
}

func (r *ChannelRepository) CreateChannel(channel model.Channel) (*model.Channel, error) {
	query := "INSERT INTO channels (name, type, user_id, server_id) VALUES ($1,$2,$3,$4) RETURNING id, created_at, updated_at"

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		channel.Name,
		channel.Type,
		channel.UserID,
		channel.ServerID,
	).Scan(
		&channel.ID,
		&channel.CreatedAt,
		&channel.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %v", err)
	}

	return &channel, nil
}
