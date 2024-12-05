package model

import (
	"time"
)

type UserRepository interface {
	CreateUserWithDefaults(user User, profile UserProfile) (*User, error)
	FindUserByField(field, value string) (*User, error)
	FindUserByFieldWithProfile(field, value string) (*User, error)

	UpdateUserProfile(profile UserProfile) (*UserProfile, error)
	DeleteUser(user User) (*User, error)
}

type UserService interface {
	RegisterUser(registerPayload UserRegisterPayload) (*User, error)
	LoginUser(u *User) (string, error)
	CheckUserCredentials(loginPayload UserLoginPayload) (*User, error)
	CheckIfEmailOrUsernameExists(registerPayload UserRegisterPayload) error
	GetUserByID(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByIDWithProfile(id string) (*User, error)

	UpdateUserProfile(userID string, userUpdatePayload UserUpdatePayload) (*UserProfile, error)
	DeleteUser(userID string, userDeletePayload UserDeletePayload) (*User, error)
}

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Profile *UserProfile `json:"profile,omitempty"`
}

type UserRegisterPayload struct {
	Username  string `json:"username" validate:"required,min=3,max=20"`
	Password  string `json:"password" validate:"required,min=6,max=130"`
	Name      string `json:"name" validate:"required"`
	ImageURL  string `json:"image_url,omitempty" validate:"omitempty,url"`
	BannerURL string `json:"banner_url,omitempty" validate:"omitempty,url"`
	Email     string `json:"email" validate:"required,email"`
}

type UserLoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ProfileStatus string

const (
	ONLINE  ProfileStatus = "ONLINE"
	BUSY    ProfileStatus = "BUSY"
	IDLE    ProfileStatus = "IDLE"
	OFFLINE ProfileStatus = "OFFLINE"
)

type UserProfile struct {
	ID        string        `json:"id"`
	UserID    string        `json:"user_id"`
	Name      string        `json:"name"`
	ImageURL  string        `json:"image_url,omitempty"`
	BannerURL string        `json:"banner_url,omitempty"`
	Bio       string        `json:"bio,omitempty"`
	Status    ProfileStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type UserUpdatePayload struct {
	Name      string        `json:"name" validate:"required"`
	ImageURL  string        `json:"image_url,omitempty"`
	BannerURL string        `json:"banner_url,omitempty"`
	Bio       string        `json:"bio,omitempty"`
	Status    ProfileStatus `json:"status" validate:"required"`
}

type UserDeletePayload struct {
	Username string `json:"username" validate:"required"`
}
