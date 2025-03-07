package servicesinfo

import "argus-backend/internal/repository/service"

type ServicesInfoInterface interface {
	GetAllServices() (*[]service.Service, error)
	AddServiceInfo(service.Service) error
	UpdateServiceInfo(service.Service) error
	DeleteService(id int) error
	GetServiceById(id int) (*service.Service, error)
}

type ServicesInfo struct {
	servicesRepository *service.ServicesRepository
}

func NewServicesInfo(servicesRepository *service.ServicesRepository) ServicesInfoInterface {
	return &ServicesInfo{
		servicesRepository: servicesRepository,
	}
}

func (si *ServicesInfo) GetAllServices() (*[]service.Service, error) {
	return si.servicesRepository.GetAllServices()
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
