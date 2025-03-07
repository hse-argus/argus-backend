package service

import "github.com/uptrace/bun"

type Service struct {
	bun.BaseModel `bun:"table:services,select:services"`
	Id            int    `bun:"id,pk,autoincrement" json:"id"`
	Name          string `bun:"name" json:"name"`
	Port          int    `bun:"port" json:"port"`
	Address       string `bun:"address" json:"address"`
}
