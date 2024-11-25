package model

import "time"

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
	FindByEmail(email string) (*User, error)
	FindByUsername(username string) (*User, error)
}

type UserService interface {
	RegisterUser(payload UserRegisterPayload) (*User, error)
	LoginUser(u *User) (string, error)
	CheckUserCredentials(loginPayload UserLoginPayload) (*User, error)
	CheckIfEmailOrUsernameExists(registerPayload UserRegisterPayload) error
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
