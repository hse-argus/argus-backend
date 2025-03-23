package app

import (
	"argus-backend/internal/logger"
	"argus-backend/internal/repository/schedule"
	"argus-backend/internal/repository/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

	path := r.URL.Path
	parts := strings.Split(path, "/")
	id, _ := strconv.Atoi(parts[2])

	serviceById, err := a.infoService.GetServiceById(id)
	if err != nil {
		http.Error(w, "Error getting service", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(serviceById)
	w.WriteHeader(http.StatusOK)
}

func (a *App) HealthCheck(w http.ResponseWriter, r *http.Request) {
	logger.Info("/health_check")

	path := r.URL.Path
	parts := strings.Split(path, "/")
	id, _ := strconv.Atoi(parts[2])

	serviceById, err := a.infoService.GetServiceById(id)
	if err != nil {
		http.Error(w, "Error getting service", 500)
		return
	}

	statusCode, err := a.clientService.SendRequest(serviceById.Address, serviceById.Port)
	if err != nil {
		http.Error(w, "Error sending request", 500)
		return
	}

	err = a.notificationService.SendWebNotification(fmt.Sprintf("Результат healthcheck: %d для сервиса %s", statusCode, serviceById.Name),
		a.connections)
	if err != nil {
		http.Error(w, "Error sending web notification", 500)
		return
	}

	err = a.notificationService.SendEmailNotification(fmt.Sprintf("Результат healthcheck: %d для сервиса %s", statusCode, serviceById.Name))
	if err != nil {
		http.Error(w, "Error sending email notification", 500)
		return
	}

	_, err = fmt.Fprintf(w, "Результат healthcheck: %d", statusCode)
	w.WriteHeader(http.StatusOK)
}

func (a *App) HandleWSConnection(w http.ResponseWriter, r *http.Request) {
	logger.Info("/ws")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("error connecting to websocket: " + err.Error())
		return
	}

	a.mu.Lock()
	a.connections[ws] = true
	a.mu.Unlock()

	defer func() {
		a.mu.Lock()
		delete(a.connections, ws)
		a.mu.Unlock()
		logger.Info("ws connection closed")
	}()

	defer ws.Close()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Info("websocket connection closed unexpectedly")
			} else {
				logger.Info("WebSocket connection closed by client")
			}
			break
		}
	}
}

func (a *App) HandleSchedule(w http.ResponseWriter, r *http.Request) {
	logger.Info("/health_check-scheduled")

	timeSchedule := schedule.Schedule{}
	data, _ := io.ReadAll(r.Body)

	err := json.Unmarshal(data, &timeSchedule)
	if err != nil {
		http.Error(w, "Something wrong with data", 400)
		return
	}

	path := r.URL.Path
	parts := strings.Split(path, "/")
	id, _ := strconv.Atoi(parts[2])

	serviceById, err := a.infoService.GetServiceById(id)
	if err != nil {
		http.Error(w, "Error getting service", 500)
		return
	}

	parsedSchedule, ok := timeSchedule.ParseScheduleDuration()
	if !ok {
		http.Error(w, "Bad time schedule", 400)
		return
	}

	err = a.StartNewJob(serviceById, parsedSchedule)
	if err != nil {
		http.Error(w, "Error starting new job", 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}
