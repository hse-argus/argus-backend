package main

import (
	"argus-backend/internal/app"
	"argus-backend/internal/config"
	"argus-backend/internal/db"
	"argus-backend/internal/server"

	"go.uber.org/fx"
)

func main() {
	addOpts := fx.Options(
		fx.Provide(config.NewConfig),
		fx.Provide(db.InitDb),
		fx.Provide(app.NewApp),
		fx.Provide(server.NewServer),
		fx.Invoke(server.RunServer),
	)
	fx.New(addOpts).Run()
}
