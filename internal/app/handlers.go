package app

import (
	"argus-backend/internal/logger"
	"argus-backend/internal/repository/service"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func (a *App) GetAllServices(w http.ResponseWriter, r *http.Request) {
	logger.Info("/get_all_services")

	services, err := a.infoService.GetAllServices()
	if err != nil {
		http.Error(w, "Error gettin all services", 500)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
}

func (a *App) UpdateService(w http.ResponseWriter, r *http.Request) {
	logger.Info("/update-service")

	updatedService := service.Service{}
	data, _ := io.ReadAll(r.Body)

	err := json.Unmarshal(data, &updatedService)
	if err != nil {
		http.Error(w, "Something wrong with input data", 400)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	err = a.infoService.UpdateServiceInfo(updatedService)
	if err != nil {
		http.Error(w, "Error updating service", 500)
	}
	w.WriteHeader(http.StatusOK)
}

func (a *App) DeleteService(w http.ResponseWriter, r *http.Request) {
	logger.Info("/delete_service")

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	err := a.infoService.DeleteService(id)
	if err != nil {
		http.Error(w, "Error deleting service", 500)
	}
	w.WriteHeader(http.StatusOK)
}

func (a *App) GetServiceById(w http.ResponseWriter, r *http.Request) {
	logger.Info("/get_service_by_id")

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	serviceById, err := a.infoService.GetServiceById(id)
	if err != nil {
		http.Error(w, "Error getting service", 500)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(serviceById)
	w.WriteHeader(http.StatusOK)
}
