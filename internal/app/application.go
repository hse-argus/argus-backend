package app

import (
	"argus-backend/internal/logger"
	servicesinfo "argus-backend/internal/service/services-info"
)

type App struct {
	infoService servicesinfo.ServicesInfoInterface
}

func NewApp(infoService servicesinfo.ServicesInfoInterface) *App {
	logger.InitLogger()
	return &App{
		infoService: infoService,
	}
}
