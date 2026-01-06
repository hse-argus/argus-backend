package user

import (
	"observer/internal/repository/service"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,select:users"`
	Id            int                `bun:"id,pk,autoincrement" json:"id"`
	Login         string             `bun:"login" json:"login"`
	Email         string             `bun:"email" json:"email"`
	Password      string             `bun:"password" json:"password"`
	Name          string             `bun:"name" json:"name"`
	Services      []*service.Service `bun:"rel:has-many,join:id=user_id"`
}
