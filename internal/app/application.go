package app

import (
	"argus-backend/internal/logger"
	clientservice "argus-backend/internal/service/client-service"
	servicesinfo "argus-backend/internal/service/services-info"
)

type App struct {
	infoService   servicesinfo.ServicesInfoInterface
	clientService *clientservice.ClientService
}

func NewApp(infoService servicesinfo.ServicesInfoInterface, service *clientservice.ClientService) *App {
	logger.InitLogger()
	return &App{
		infoService:   infoService,
		clientService: service,
	}
}
