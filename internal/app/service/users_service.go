package service

import (
	"fmt"

	"github.com/razaq-himawan/chat-app-api/internal/app/model"
	"github.com/razaq-himawan/chat-app-api/internal/auth"
)

type UserService struct {
	userRepo model.UserRepository
}

func NewUserService(userRepo model.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) RegisterUser(registerPayload model.UserRegisterPayload) (*model.User, error) {
	hashedPassword, err := auth.HashPassword(registerPayload.Password)
	if err != nil {
		return nil, err
	}

	// TODO imageURL thing

	createdUser, err := s.userRepo.CreateUser(model.User{
		Username: registerPayload.Username,
		Password: hashedPassword,
		Name:     registerPayload.Name,
		ImageURL: registerPayload.ImageURL,
		Email:    registerPayload.Email,
	})
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *UserService) LoginUser(u *model.User) (string, error) {

	token, err := auth.CreateJWT(u.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) CheckUserCredentials(loginPayload model.UserLoginPayload) (*model.User, error) {
	u, err := s.GetUserByEmail(loginPayload.Email)
	if err != nil {
		return nil, fmt.Errorf("email or password do not match")
	}

	if !auth.ComparePasswords(u.Password, []byte(loginPayload.Password)) {
		return nil, fmt.Errorf("email or password do not match")
	}

	return u, nil
}

func (s *UserService) CheckIfEmailOrUsernameExists(registerPayload model.UserRegisterPayload) error {
	_, err := s.userRepo.FindByEmail(registerPayload.Email)
	if err == nil {
		return fmt.Errorf("user with email %v already exists", registerPayload.Email)
	}

	_, err = s.userRepo.FindByUsername(registerPayload.Username)
	if err == nil {
		return fmt.Errorf("user with username %v already exists", registerPayload.Username)
	}

	return nil
}

func (s *UserService) GetUserByID(id string) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	return s.userRepo.FindByEmail(email)
}

func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	return s.userRepo.FindByUsername(username)
}
