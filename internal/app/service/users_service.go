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

	createdUser, err := s.userRepo.CreateUserWithDefaults(
		model.User{
			Username: registerPayload.Username,
			Password: hashedPassword,
			Email:    registerPayload.Email,
		},
		model.UserProfile{
			Name:      registerPayload.Name,
			ImageURL:  registerPayload.ImageURL,
			BannerURL: registerPayload.BannerURL,
			Status:    model.OFFLINE,
		},
	)
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
	_, err := s.GetUserByEmail(registerPayload.Email)
	if err == nil {
		return fmt.Errorf("user with email %v already exists", registerPayload.Email)
	}

	_, err = s.GetUserByUsername(registerPayload.Username)
	if err == nil {
		return fmt.Errorf("user with username %v already exists", registerPayload.Username)
	}

	return nil
}

func (s *UserService) GetUserByID(id string) (*model.User, error) {
	return s.userRepo.FindUserByField("id", id)
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	return s.userRepo.FindUserByField("email", email)
}

func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	return s.userRepo.FindUserByFieldWithProfile("username", username)
}

func (s *UserService) GetUserByIDWithProfile(id string) (*model.User, error) {
	return s.userRepo.FindUserByFieldWithProfile("id", id)
}

func (s *UserService) UpdateUserProfile(userID string, userUpdatePayload model.UserUpdatePayload) (*model.UserProfile, error) {
	if !s.isStatusValid(string(userUpdatePayload.Status)) {
		return nil, fmt.Errorf("invalid status")
	}

	return s.userRepo.UpdateUserProfile(model.UserProfile{
		Name:      userUpdatePayload.Name,
		ImageURL:  userUpdatePayload.ImageURL,
		BannerURL: userUpdatePayload.BannerURL,
		Bio:       userUpdatePayload.Bio,
		Status:    userUpdatePayload.Status,
		UserID:    userID,
	})
}

func (s *UserService) DeleteUser(userID string, userDeletePayload model.UserDeletePayload) (*model.User, error) {
	us, err := s.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if userDeletePayload.Username != us.Username {
		return nil, fmt.Errorf("username do not match")
	}

	u, err := s.userRepo.DeleteUser(model.User{
		ID:       userID,
		Username: userDeletePayload.Username,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to delete user")
	}

	return u, nil
}

func (s *UserService) isStatusValid(status string) bool {
	var validStatus = []string{
		string(model.ONLINE),
		string(model.BUSY),
		string(model.IDLE),
		string(model.OFFLINE),
	}

	for _, s := range validStatus {
		if s == status {
			return true
		}
	}
	return false
}
