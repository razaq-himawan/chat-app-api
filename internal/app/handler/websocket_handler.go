package handler

// TODO: MAKE THIS HANDLER DYNAMIC

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	ws "github.com/coder/websocket"
	"github.com/razaq-himawan/chat-app-api/internal/app/model"
	"github.com/razaq-himawan/chat-app-api/internal/auth"
	"github.com/razaq-himawan/chat-app-api/internal/websocket"
	"github.com/razaq-himawan/chat-app-api/utils"
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenString := utils.GetTokenFromCookie(r)

	userID, err := auth.GetUserIDFromToken(tokenString)
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}

	conn, err := ws.Accept(w, r, nil)
	if err != nil {
		log.Println("Failed to accept WebSocket connection:", err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to open WebSocket connection"))
		return
	}

	wsServer := websocket.GetWebSocketServer()

	go wsServer.Start(ctx)

	client := &model.WebSocketUser{
		UserID:         userID,
		Conn:           conn,
		Type:           model.DM,
		ConversationID: "conv1",
		IsOnline:       true,
	}
	wsServer.Register <- client

	defer func() {
		wsServer.Unregister <- client
		if err := conn.Close(ws.StatusNormalClosure, "Connection closed"); err != nil {
			log.Println("Error closing WebSocket connection:", err)
		}
	}()

	for {
		_, messageBytes, err := conn.Read(ctx)
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var message model.Message
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			log.Println("Error unmarshaling message:", err)
			continue
		}

		wsServer.Broadcast <- &model.Message{
			ID:             "test",
			Content:        message.Content,
			MemberID:       client.UserID,
			ConversationID: message.ConversationID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
	}
}
