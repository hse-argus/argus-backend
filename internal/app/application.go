package app

import (
	"argus-backend/internal/logger"
	clientservice "argus-backend/internal/service/client-service"
	notificationservice "argus-backend/internal/service/notification-service"
	servicesinfo "argus-backend/internal/service/services-info"
	"github.com/gorilla/websocket"
	"sync"
)

type App struct {
	infoService         servicesinfo.ServicesInfoInterface
	clientService       *clientservice.ClientService
	notificationService *notificationservice.WebNotificationService

	connections map[*websocket.Conn]bool
	mu          *sync.RWMutex
}

func NewApp(infoService servicesinfo.ServicesInfoInterface, service *clientservice.ClientService) *App {
	logger.InitLogger()
	return &App{
		infoService:   infoService,
		clientService: service,
		connections:   make(map[*websocket.Conn]bool),
		mu:            &sync.RWMutex{},
	}
}
