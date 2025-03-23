package notificationservice

import (
	"argus-backend/internal/logger"
	"github.com/gorilla/websocket"
)

type NotificationService struct{}

func NewWebNotificationService() *NotificationService {
	return &NotificationService{}
}

func (wns *NotificationService) SendWebNotification(text string, connections map[*websocket.Conn]bool) error {
	for connection := range connections {
		err := connection.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			logger.Error("error sending ws event: " + err.Error())
			return err
		}
	}
	return nil
}
