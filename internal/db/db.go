package db

import (
	"observer/internal/config"
	"observer/internal/logger"
	"observer/internal/repository/service"
	"observer/internal/repository/user"
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/fx"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/uptrace/bun"
)

func InitDb(lc fx.Lifecycle, config *config.Config) *bun.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable",
		config.PostgresUser, config.PostgresPassword, config.PostgresPort, config.PostgresDb)
	sqldb, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Error(fmt.Sprintf("error connecting to database: %v", err))
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	_, err = db.NewCreateTable().
		IfNotExists().
		Model((*service.Service)(nil)).
		Exec(context.Background())

	if err != nil {
		logger.Error(fmt.Sprintf("error creating table service: %v", err))
	}

	_, err = db.NewCreateTable().
		IfNotExists().
		Model((*user.User)(nil)).
		Exec(context.Background())

	if err != nil {
		logger.Error(fmt.Sprintf("error creating table user: %v", err))
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err = db.Close()
			if err != nil {
				return err
			}

			err = sqldb.Close()
			if err != nil {
				return err
			}

			return nil
		},
	})

	return db
}
