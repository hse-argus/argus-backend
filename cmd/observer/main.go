package main

import (
	_ "observer/docs"
	"observer/internal/app"
	"observer/internal/config"
	"observer/internal/db"
	"observer/internal/logger"
	"observer/internal/repository/service"
	"observer/internal/repository/user"
	"observer/internal/server"
	clientservice "observer/internal/service/client-service"
	"observer/internal/service/notification-service"
	servicesinfo "observer/internal/service/services-info"
	userservice "observer/internal/service/user-service"

	"go.uber.org/fx"
)

// @title Swagger observer
// @version 1.0
// @host localhost:8080
// @BasePath /
func main() {
	addOpts := fx.Options(
		fx.Provide(config.NewConfig,
			clientservice.NewClientService,
			db.InitDb,
			service.NewServicesRepository,
			user.NewUserRepository,
			userservice.NewUserService,
			servicesinfo.NewServicesInfo,
			notificationservice.NewWebNotificationService,
			app.NewApp,
			server.NewServer),
		fx.Invoke(logger.InitLogger, app.InvokeScheduler, server.RunServer),
	)
	fx.New(addOpts).Run()
}
