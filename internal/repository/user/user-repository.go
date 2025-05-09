package user

import (
	customerrors "argus-backend/internal/errors"
	"argus-backend/internal/logger"
	"context"
	"fmt"
	"github.com/uptrace/bun"
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) GetUserByLogin(login string) (*User, error) {
	var user User
	err := ur.db.NewSelect().
		Model((*User)(nil)).
		Where("login = ?", login).
		Scan(context.Background(), &user)
	if err != nil {
		logger.Info(fmt.Sprintf("get user error: %v", err))
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) AddUser(user User) (int, error) {
	existing, _ := ur.GetUserByLogin(user.Login)
	if existing != nil {
		logger.Info("user already exists")
		return 0, customerrors.AlreadyExistsError{}
	}
	var id int

	err := ur.db.NewInsert().
		Model(&user).
		Returning("id").
		Scan(context.Background(), &id)
	if err != nil {
		logger.Info(fmt.Sprintf("add user error: %v", err))
		return 0, err
	}

	return id, nil
}

func (ur *UserRepository) GetUser(login string, password string) (*User, error) {
	var user User
	err := ur.db.NewSelect().
		Model((*User)(nil)).
		Where("login = ? and password = ?", login, password).
		Scan(context.Background(), &user)
	if err != nil {
		logger.Info(fmt.Sprintf("get user error: %v", err))
		return &User{}, err
	}

	return &user, nil
}
