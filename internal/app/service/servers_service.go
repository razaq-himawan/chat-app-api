package service

import (
	"database/sql"
	"log"

	"github.com/razaq-himawan/chat-app-api/internal/app/model"
)

type ServerService struct {
	db          *sql.DB
	serverRepo  model.ServerRepository
	memberRepo  model.MemberRepository
	channelRepo model.ChannelRepository
}

func NewServerService(
	db *sql.DB,
	serverRepo model.ServerRepository,
	memberRepo model.MemberRepository,
	channelRepo model.ChannelRepository,
) *ServerService {
	return &ServerService{
		db:          db,
		serverRepo:  serverRepo,
		memberRepo:  memberRepo,
		channelRepo: channelRepo,
	}
}

func (s *ServerService) CreateServerWithMembersAndChannels(createServerPayload model.CreateServerPayload, userID string) (*model.ServerModel, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			log.Printf("recovered on panic: %v", p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	server, err := s.serverRepo.CreateServerTx(tx, model.ServerModel{
		Name:   createServerPayload.Name,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	_, err = s.memberRepo.CreateMemberTx(tx, model.Member{
		Role:     model.ADMIN,
		UserID:   userID,
		ServerID: server.ID,
	})
	if err != nil {
		return nil, err
	}

	_, err = s.channelRepo.CreateChannelTx(tx, model.Channel{
		Name:     "general",
		Type:     model.TEXT,
		UserID:   userID,
		ServerID: server.ID,
	})
	if err != nil {
		return nil, err
	}

	return server, nil
}
