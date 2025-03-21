package notificationservice

import (
	"argus-backend/internal/logger"
	"github.com/gorilla/websocket"
)

type WebNotificationService struct{}

func NewWebNotificationService() *WebNotificationService {
	return &WebNotificationService{}
}

func (wns *WebNotificationService) SendNotification(text string, connections map[*websocket.Conn]bool) error {
	for connection := range connections {
		err := connection.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			logger.Error("error sending ws event: " + err.Error())
			return err
		}
	}
	return nil
}
