package repository

import (
	"database/sql"
	"fmt"

	"github.com/razaq-himawan/chat-app-api/internal/app/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user model.User) (*model.User, error) {
	stmt, err := r.db.Prepare("INSERT INTO users (username, password, name, image_url, email) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at, updated_at")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		user.Username,
		user.Password,
		user.Name,
		user.ImageURL,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %v", err)
	}

	return &user, nil
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	return r.findUserByField("id", id)
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	return r.findUserByField("email", email)
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	return r.findUserByField("username", username)
}

func (r *UserRepository) findUserByField(field, value string) (*model.User, error) {
	query := fmt.Sprintf("SELECT * FROM users WHERE %s = $1", field)
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(value)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("user not found")
	}

	u, err := scanRowIntoUser(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan user: %v", err)
	}

	return u, nil
}

func scanRowIntoUser(rows *sql.Rows) (*model.User, error) {
	user := new(model.User)

	err := rows.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Name,
		&user.ImageURL,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
