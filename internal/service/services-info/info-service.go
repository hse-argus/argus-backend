package servicesinfo

import "argus-backend/internal/repository/service"

type ServicesInfoInterface interface {
	GetAllServices() (*[]service.Service, error)
	AddServiceInfo(service.Service) error
	// UpdateServiceInfo(service.Service) error
	// DeleteService() error
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
