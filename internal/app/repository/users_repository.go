package repository

import (
	"database/sql"
	"fmt"

	"github.com/razaq-himawan/chat-app-api/internal/app/model"
	"github.com/razaq-himawan/chat-app-api/internal/app/repository/helper"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUserWithDefaults(user model.User, profile model.UserProfile) (*model.User, error) {
	result, err := helper.ExecWithTx(r.db, func(tx *sql.Tx) (*model.User, error) {
		userQuery := `
			INSERT INTO users (username, password, email) 
			VALUES ($1, $2, $3) 
			RETURNING id, created_at, updated_at
		`
		err := tx.QueryRow(
			userQuery,
			user.Username,
			user.Password,
			user.Email,
		).Scan(
			&user.ID,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		profileQuery := `
			INSERT INTO profiles (user_id, name, image_url, banner_url, bio, status) 
			VALUES ($1, $2, $3, $4, $5, $6) 
			RETURNING id, user_id, created_at, updated_at
		`
		err = tx.QueryRow(
			profileQuery,
			user.ID,
			profile.Name,
			profile.ImageURL,
			profile.BannerURL,
			profile.Bio,
			profile.Status,
		).Scan(
			&profile.ID,
			&profile.UserID,
			&profile.CreatedAt,
			&profile.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create profile: %w", err)
		}

		user.Profile = &profile

		return &user, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user with profile: %w", err)
	}

	return result, nil
}

func (r *UserRepository) FindUserByField(field, value string) (*model.User, error) {
	query := fmt.Sprintf("SELECT id, username, password, email, created_at, updated_at FROM users WHERE %s = $1", field)

	user := &model.User{}

	err := r.db.QueryRow(query, value).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to scan user: %v", err)
	}

	return user, nil
}

func (r *UserRepository) FindUserByFieldWithProfile(field, value string) (*model.User, error) {
	query := fmt.Sprintf(`
		SELECT 
			u.id, u.username, u.email, u.created_at, u.updated_at,
			p.id, p.user_id, p.name, p.image_url, p.banner_url, p.bio, p.status, p.created_at, p.updated_at
		FROM users u
		LEFT JOIN profiles p ON u.id = p.user_id
		WHERE u.%s = $1
	`, field)

	user := &model.User{}
	profile := &model.UserProfile{}
	err := r.db.QueryRow(query, value).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&profile.ID,
		&profile.UserID,
		&profile.Name,
		&profile.ImageURL,
		&profile.BannerURL,
		&profile.Bio,
		&profile.Status,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found")
		}
		return nil, fmt.Errorf("failed to fetch user with profile: %v", err)
	}

	if profile.ID != "" {
		user.Profile = profile
	}

	return user, nil
}

// TODO: Update User

func (r *UserRepository) UpdateUserProfile(profile model.UserProfile) (*model.UserProfile, error) {
	query := `
		UPDATE profiles
		SET name = $1, image_url = $2, banner_url = $3, bio = $4, status = $5
		WHERE user_id = $6
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		profile.Name,
		profile.ImageURL,
		profile.BannerURL,
		profile.Bio,
		profile.Status,
		profile.UserID,
	).Scan(
		&profile.ID,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found")
		}
		return nil, fmt.Errorf("failed to fetch user with profile: %v", err)
	}

	return &profile, nil
}
