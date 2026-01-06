package notificationservice

import (
	"observer/internal/logger"
	"github.com/gorilla/websocket"
)

type NotificationService struct{}

func NewWebNotificationService() *NotificationService {
	return &NotificationService{}
}

func (wns *NotificationService) SendWebNotification(text string,
	connections map[string]*websocket.Conn,
	userLogin string) error {
	connection := connections[userLogin]
	if connection == nil {
		return nil
	}
	err := connections[userLogin].WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logger.Error("error sending ws event: " + err.Error())
		return err
	}
	logger.Info("sent web event: " + text)

	return nil
}
