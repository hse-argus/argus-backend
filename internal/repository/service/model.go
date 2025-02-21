package service

import "github.com/uptrace/bun"

type Service struct {
	bun.BaseModel `bun:"table:services,select:services"`
	Id            int    `bun:"id,pk,autoincrement" json:"id"`
	Name          string `bun:"name" json:"Name"`
	Port          int    `bun:"port" json:"Port"`
}
