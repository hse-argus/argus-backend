package app

import (
	"argus-backend/internal/config"
	"argus-backend/internal/db"
)

type App struct {
	conf *config.Config
	db *db.Db
}

func NewApp(conf *config.Config, database *db.Db) *App {
	return &App{
		conf: conf,
		db: database,
	}
}