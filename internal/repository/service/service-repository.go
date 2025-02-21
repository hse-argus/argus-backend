package service

import (
	"context"

	"github.com/uptrace/bun"
)

type ServicesRepository struct {
	db *bun.DB
}

func NewServicesRepository(db *bun.DB) *ServicesRepository {
	return &ServicesRepository{
		db: db,
	}
}

func (sr *ServicesRepository) GetAllServices() (*[]Service, error) {
	service := make([]Service, 0)

	err := sr.db.NewSelect().
		Model(&service).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (sr *ServicesRepository) AddServiceInfo(newService Service) error {
	_, err := sr.db.NewInsert().
		Model(&newService).
		Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}
