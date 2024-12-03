package repository

import (
	"database/sql"
	"fmt"

	"github.com/razaq-himawan/chat-app-api/internal/app/model"
)

type MemberRepository struct {
	db *sql.DB
}

func NewMemberRepository(db *sql.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

func (r *MemberRepository) CreateMember(member model.Member) (*model.Member, error) {
	query := "INSERT INTO members (role, user_id, server_id) VALUES ($1,$2,$3) RETURNING id, created_at, updated_at"

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		member.Role,
		member.UserID,
		member.ServerID,
	).Scan(
		&member.ID,
		&member.CreatedAt,
		&member.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %v", err)
	}

	return &member, nil
}
