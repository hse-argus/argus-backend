package main

import (
	"argus-backend/internal/app"
	"argus-backend/internal/config"
	"argus-backend/internal/db"
	"argus-backend/internal/repository/service"
	"argus-backend/internal/server"
	servicesinfo "argus-backend/internal/service/services-info"

	"go.uber.org/fx"
)

func main() {
	addOpts := fx.Options(
		fx.Provide(config.NewConfig),
		fx.Provide(db.InitDb),
		fx.Provide(service.NewServicesRepository),
		fx.Provide(servicesinfo.NewServicesInfo),
		fx.Provide(app.NewApp),
		fx.Provide(server.NewServer),
		fx.Invoke(server.RunServer),
	)
	fx.New(addOpts).Run()
}
