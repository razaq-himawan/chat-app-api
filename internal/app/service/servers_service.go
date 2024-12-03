package service

import (
	"github.com/razaq-himawan/chat-app-api/internal/app/model"
)

type ServerService struct {
	serverRepo model.ServerRepository
}

func NewServerService(serverRepo model.ServerRepository) *ServerService {
	return &ServerService{serverRepo: serverRepo}
}

func (s *ServerService) CreateServerWithMembersAndChannels(createServerPayload model.CreateServerPayload, userID string) (*model.ServerModel, error) {
	server, err := s.serverRepo.CreateServerWithDefaults(model.ServerModel{
		Name:   createServerPayload.Name,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}
	return server, nil
}
