package app

import (
	"argus-backend/internal/logger"
	"argus-backend/internal/repository/service"
	"encoding/json"
	"io"
	"net/http"
)

func (a *App) GetAllServices(w http.ResponseWriter, r *http.Request) {
	logger.Info("/get_all_services")

	services, err := a.infoService.GetAllServices()
	if err != nil {
		http.Error(w, "Error gettin all services", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(services)
	w.WriteHeader(http.StatusOK)
}

func (a *App) AddService(w http.ResponseWriter, r *http.Request) {
	logger.Info("/add_service")

	newService := service.Service{}
	data, _ := io.ReadAll(r.Body)

	err := json.Unmarshal(data, &newService)
	if err != nil {
		http.Error(w, "Something wrong with data", 400)
		return
	}

	err = a.infoService.AddServiceInfo(newService)
	if err != nil {
		http.Error(w, "Error adding service", 500)
	}
	w.WriteHeader(http.StatusOK)
}
