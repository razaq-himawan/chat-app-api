package repository

import (
	"database/sql"
	"fmt"

	"github.com/razaq-himawan/chat-app-api/internal/app/model"
)

type ServerRepository struct {
	db *sql.DB
}

func NewServerRepository(db *sql.DB) *ServerRepository {
	return &ServerRepository{db: db}
}

func (r *ServerRepository) CreateServerTx(tx *sql.Tx, server model.ServerModel) (*model.ServerModel, error) {
	stmt, err := tx.Prepare("INSERT INTO servers (name, user_id) VALUES ($1,$2) RETURNING id, invite_code, created_at, updated_at")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		server.Name,
		server.UserID,
	).Scan(
		&server.ID,
		&server.InviteCode,
		&server.CreatedAt,
		&server.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %v", err)
	}

	return &server, nil
}
