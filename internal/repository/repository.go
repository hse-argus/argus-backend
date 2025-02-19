package repository

import "github.com/uptrace/bun"

type Service struct {
	bun.BaseModel `bun:"table:services,select:services"`
	Id   int    `bun:"id,pk,autoincrement"`
	Name string `bun:"name"`
	Port int    `bun:"port"`
}
