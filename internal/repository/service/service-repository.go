package service

import (
	"observer/internal/logger"
	"context"
	"github.com/google/uuid"

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

func (sr *ServicesRepository) GetAllServices(userId int) (*[]Service, error) {
	service := make([]Service, 0)

	err := sr.db.NewSelect().
		Model(&service).
		Where("user_id = ?", userId).
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

func (sr *ServicesRepository) UpdateServiceInfo(newService Service) error {
	_, err := sr.db.NewUpdate().
		Model(&newService).
		Where("id = ?", newService.Id).
		OmitZero().
		Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (sr *ServicesRepository) DeleteServiceInfo(id int) error {
	_, err := sr.db.NewDelete().
		Model((*Service)(nil)).
		Where("id = ?", id).
		Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (sr *ServicesRepository) GetServiceById(id int) (*Service, error) {
	service := new(Service)
	err := sr.db.NewSelect().
		Model(service).
		Where("id = ?", id).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (sr *ServicesRepository) AddJobID(id int, jobID uuid.UUID) error {
	_, err := sr.db.NewUpdate().
		Model((*Service)(nil)).
		Where("id = ?", id).
		Set("job_id = ?", jobID).
		Exec(context.Background())
	if err != nil {
		logger.Error("error updating job_id: " + err.Error())
		return err
	}

	return nil
}

func (sr *ServicesRepository) getJobID(id int) (uuid.UUID, error) {
	var service Service

	err := sr.db.NewSelect().
		Model(&service).
		Where("id = ?", id).
		Scan(context.Background())
	if err != nil {
		logger.Info("error getting job_id: " + err.Error())
		return uuid.Nil, err
	}

	return service.JobID, nil
}

func (sr *ServicesRepository) DeleteJobID(id int) (uuid.UUID, error) {
	jobID, err := sr.getJobID(id)
	if err != nil {
		return uuid.Nil, err
	}

	_, err = sr.db.NewUpdate().
		Model((*Service)(nil)).
		Where("id = ?", id).
		Set("job_id = ?", uuid.Nil).
		Exec(context.Background())

	logger.Info("deleted job_id: " + jobID.String())
	if err != nil {
		logger.Error("error deleting job_id: " + err.Error())
		return uuid.Nil, err
	}

	return jobID, nil
}
