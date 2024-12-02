package websocket

import (
	"context"
	"log"
	"sync"

	"github.com/coder/websocket"
	"github.com/razaq-himawan/chat-app-api/internal/app/model"
)

type WebSocketServer struct {
	DmClients      map[string]*model.WebSocketUser
	ChannelClients map[string]*model.WebSocketUser
	Broadcast      chan *model.Message
	Register       chan *model.WebSocketUser
	Unregister     chan *model.WebSocketUser
	mu             sync.RWMutex
}

var wsServer *WebSocketServer
var once sync.Once

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		DmClients:      make(map[string]*model.WebSocketUser),
		ChannelClients: make(map[string]*model.WebSocketUser),
		Broadcast:      make(chan *model.Message, 100),
		Register:       make(chan *model.WebSocketUser, 100),
		Unregister:     make(chan *model.WebSocketUser, 100),
	}
}

func (s *WebSocketServer) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("WebSocket server shutting down...")
			s.disconnectAllClients()
			return
		case client := <-s.Register:
			s.registerClient(client)

		case client := <-s.Unregister:
			s.unregisterClient(client)

		case message := <-s.Broadcast:
			s.handleBroadcastMessage(ctx, message)
		}
	}
}

func (s *WebSocketServer) registerClient(client *model.WebSocketUser) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if client.Type != model.DM && client.Type != model.CHANNEL {
		log.Printf("Invalid client type: %s", client.Type)
		return
	}

	client.IsOnline = true
	if client.Type == model.DM {
		s.DmClients[client.UserID] = client
	} else if client.Type == model.CHANNEL {
		s.ChannelClients[client.UserID] = client
	}
}

func (s *WebSocketServer) unregisterClient(client *model.WebSocketUser) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if client.Type == model.DM {
		delete(s.DmClients, client.UserID)
	} else if client.Type == model.CHANNEL {
		delete(s.ChannelClients, client.UserID)
	}
	client.IsOnline = false
}

func (server *WebSocketServer) handleBroadcastMessage(ctx context.Context, message *model.Message) {
	if message.ConversationID != "" {
		server.SendMessageToConversation(ctx, message.ConversationID, message.Content)
	} else if message.ChannelID != "" {
		server.SendMessageToChannel(ctx, message.ChannelID, message.Content)
	} else {
		log.Println("Invalid message: neither ConversationID nor ChannelID provided")
	}
}

func (server *WebSocketServer) SendMessageToConversation(ctx context.Context, conversationID, message string) {
	server.mu.RLock()
	defer server.mu.RUnlock()

	for _, client := range server.DmClients {
		if client.ConversationID == conversationID && client.IsOnline {
			if err := client.Conn.Write(ctx, websocket.MessageText, []byte(message)); err != nil {
				server.handleConnectionError(client, client.UserID)
			}
		}
	}
}

func (server *WebSocketServer) SendMessageToChannel(ctx context.Context, channelID, message string) {
	server.mu.RLock()
	defer server.mu.RUnlock()

	for _, client := range server.ChannelClients {
		if client.ChannelID == channelID && client.IsOnline {
			if err := client.Conn.Write(ctx, websocket.MessageText, []byte(message)); err != nil {
				server.handleConnectionError(client, client.UserID)
			}
		}
	}
}

func (server *WebSocketServer) handleConnectionError(client *model.WebSocketUser, userID string) {
	log.Printf("User %s is offline, removing from active clients", userID)

	client.IsOnline = false
	safeClose(client.Conn, userID)

	server.mu.Lock()
	defer server.mu.Unlock()

	if client.Type == model.DM {
		delete(server.DmClients, userID)
	} else if client.Type == model.CHANNEL {
		delete(server.ChannelClients, userID)
	}
}

func (s *WebSocketServer) disconnectAllClients() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for userID, client := range s.DmClients {
		safeClose(client.Conn, userID)
		delete(s.DmClients, userID)
	}
	for userID, client := range s.ChannelClients {
		safeClose(client.Conn, userID)
		delete(s.ChannelClients, userID)
	}
	log.Println("All clients disconnected.")
}

func safeClose(conn *websocket.Conn, userID string) {
	if conn == nil {
		return
	}
	if err := conn.Close(websocket.StatusNormalClosure, "closing connection due to error"); err != nil {
		log.Printf("Error closing WebSocket connection for user %s: %v", userID, err)
	}
}

func GetWebSocketServer() *WebSocketServer {
	once.Do(func() {
		wsServer = NewWebSocketServer()
	})
	return wsServer
}
