package main

import (
	_ "argus-backend/docs"
	"argus-backend/internal/app"
	"argus-backend/internal/config"
	"argus-backend/internal/db"
	"argus-backend/internal/logger"
	"argus-backend/internal/repository/service"
	"argus-backend/internal/server"
	clientservice "argus-backend/internal/service/client-service"
	"argus-backend/internal/service/notification-service"
	servicesinfo "argus-backend/internal/service/services-info"

	"go.uber.org/fx"
)

// @title Swagger Argus-Backend
// @version 1.0
// @host localhost:8080
// @BasePath /
func main() {
	addOpts := fx.Options(
		fx.Provide(config.NewConfig,
			clientservice.NewClientService,
			db.InitDb,
			service.NewServicesRepository,
			servicesinfo.NewServicesInfo,
			notificationservice.NewWebNotificationService,
			app.NewApp,
			server.NewServer),
		fx.Invoke(logger.InitLogger, app.InvokeScheduler, server.RunServer),
	)
	fx.New(addOpts).Run()
}
