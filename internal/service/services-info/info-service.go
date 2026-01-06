package servicesinfo

import (
	"observer/internal/repository/service"
	"github.com/google/uuid"
)

type ServicesInfoInterface interface {
	GetAllServices(userId int) (*[]service.Service, error)
	AddServiceInfo(service.Service) error
	UpdateServiceInfo(service.Service) error
	DeleteService(id int) error
	GetServiceById(id int) (*service.Service, error)
	AddJob(id int, jobID uuid.UUID) error
	DeleteJob(id int) (uuid.UUID, error)
}

type ServicesInfo struct {
	servicesRepository *service.ServicesRepository
}

func NewServicesInfo(servicesRepository *service.ServicesRepository) ServicesInfoInterface {
	return &ServicesInfo{
		servicesRepository: servicesRepository,
	}
}

func (si *ServicesInfo) GetAllServices(userId int) (*[]service.Service, error) {
	return si.servicesRepository.GetAllServices(userId)
}

func (si *ServicesInfo) AddServiceInfo(newService service.Service) error {
	return si.servicesRepository.AddServiceInfo(newService)
}

func (si *ServicesInfo) UpdateServiceInfo(newService service.Service) error {
	return si.servicesRepository.UpdateServiceInfo(newService)
}

func (si *ServicesInfo) DeleteService(id int) error {
	return si.servicesRepository.DeleteServiceInfo(id)
}

func (si *ServicesInfo) GetServiceById(id int) (*service.Service, error) {
	return si.servicesRepository.GetServiceById(id)
}

func (si *ServicesInfo) AddJob(id int, jobID uuid.UUID) error {
	return si.servicesRepository.AddJobID(id, jobID)
}

func (si *ServicesInfo) DeleteJob(id int) (uuid.UUID, error) {
	return si.servicesRepository.DeleteJobID(id)
}
