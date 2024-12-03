package repository

import (
	"database/sql"
	"fmt"

	"github.com/razaq-himawan/chat-app-api/internal/app/model"
	"github.com/razaq-himawan/chat-app-api/internal/app/repository/helper"
)

type ServerRepository struct {
	db *sql.DB
}

func NewServerRepository(db *sql.DB) *ServerRepository {
	return &ServerRepository{db: db}
}

func (r *ServerRepository) CreateServerWithDefaults(server model.ServerModel) (*model.ServerModel, error) {
	result, err := helper.ExecWithTx(r.db, func(tx *sql.Tx) (*model.ServerModel, error) {
		serverQuery := "INSERT INTO servers (name, user_id) VALUES ($1,$2) RETURNING id, invite_code, created_at, updated_at"
		err := tx.QueryRow(
			serverQuery,
			server.Name,
			server.UserID,
		).Scan(
			&server.ID,
			&server.InviteCode,
			&server.CreatedAt,
			&server.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create server: %v", err)
		}

		memberQuery := "INSERT INTO members (role, user_id, server_id) VALUES ($1,$2,$3) RETURNING *"
		var member model.Member
		err = tx.QueryRow(
			memberQuery,
			model.ADMIN,
			server.UserID,
			server.ID,
		).Scan(
			&member.ID,
			&member.Role,
			&member.UserID,
			&member.ServerID,
			&member.CreatedAt,
			&member.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create member: %v", err)
		}

		channelQuery := "INSERT INTO channels (name, type, user_id, server_id) VALUES ($1,$2,$3,$4) RETURNING *"
		var channel model.Channel
		err = tx.QueryRow(
			channelQuery,
			"general",
			model.TEXT,
			server.UserID,
			server.ID,
		).Scan(
			&channel.ID,
			&channel.Name,
			&channel.Type,
			&channel.UserID,
			&channel.ServerID,
			&channel.CreatedAt,
			&channel.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create channel: %v", err)
		}

		server.Members = []model.Member{member}
		server.Channel = []model.Channel{channel}

		return &server, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create server with defaults: %v", err)
	}

	return result, nil
}
