package db

import (
	"argus-backend/internal/config"
	"argus-backend/internal/repository"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/uptrace/bun"
)

type Db struct {
	db *bun.DB
}

func InitDb(config *config.Config) *Db {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable",
		config.PostgresUser, config.PostgresPassword, config.PostgresPort, config.PostgresDb)
	sqldb, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	_, err = db.NewCreateTable().
		IfNotExists().
		Model((*repository.Service)(nil)).
		Exec(context.Background())

	if err != nil {
		fmt.Printf("error: %v", err)
	}
	return &Db{db: db}
}
