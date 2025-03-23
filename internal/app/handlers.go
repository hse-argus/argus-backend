package app

import (
	"argus-backend/internal/logger"
	"argus-backend/internal/repository/schedule"
	"argus-backend/internal/repository/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (a *App) GetAllServices(c *gin.Context) {
	logger.Info("/get_all_services")

	services, err := a.infoService.GetAllServices()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("error getting all services"))
		return
	}

	c.JSON(http.StatusOK, services)
}

func (a *App) AddService(c *gin.Context) {
	logger.Info("/add_service")

	newService := service.Service{}
	if err := c.BindJSON(&newService); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("error parsing body"))
		return
	}

	err := a.infoService.AddServiceInfo(newService)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("error adding service"))
		return
	}
	c.Status(http.StatusOK)
}

func (a *App) UpdateService(c *gin.Context) {
	logger.Info("/update-service")

	updatedService := service.Service{}
	if err := c.BindJSON(&updatedService); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("error parsing body"))
	}

	err := a.infoService.UpdateServiceInfo(updatedService)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("error updating service"))
	}
	c.Status(http.StatusOK)
}

func (a *App) DeleteService(c *gin.Context) {
	logger.Info("/delete_service")

	id, _ := strconv.Atoi(c.Query("id"))
	err := a.infoService.DeleteService(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("error deleting service"))
	}
	c.Status(http.StatusOK)
}

func (a *App) GetServiceById(c *gin.Context) {
	logger.Info("/get_service_by_id")

	id, _ := strconv.Atoi(c.Param("id"))
	serviceById, err := a.infoService.GetServiceById(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("error getting service"))
		return
	}

	c.JSON(http.StatusOK, serviceById)
}

func (a *App) HealthCheck(c *gin.Context) {
	logger.Info("/health_check")

	id, _ := strconv.Atoi(c.Param("id"))
	serviceById, err := a.infoService.GetServiceById(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("error getting service"))
		return
	}

	statusCode, err := a.clientService.SendRequest(serviceById.Address, serviceById.Port)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("error sending request"))
		return
	}

	err = a.notificationService.SendWebNotification(fmt.Sprintf("Результат healthcheck: %d для сервиса %s",
		statusCode,
		serviceById.Name),
		a.connections)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("error sending web notification"))
		return
	}

	err = a.notificationService.SendEmailNotification(fmt.Sprintf("Результат healthcheck: %d для сервиса %s",
		statusCode,
		serviceById.Name))
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("error sending email notification"))
		return
	}

	c.Status(http.StatusOK)
}

func (a *App) HandleWSConnection(c *gin.Context) {
	logger.Info("/ws")
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
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

func (a *App) HandleSchedule(c *gin.Context) {
	logger.Info("/health_check-scheduled")

	timeSchedule := schedule.Schedule{}
	if err := c.BindJSON(&timeSchedule); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("error parsing body"))
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	serviceById, err := a.infoService.GetServiceById(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("error getting service"))
		return
	}

	parsedSchedule, ok := timeSchedule.ParseScheduleDuration()
	if !ok {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("unexisted schedule"))
		return
	}

	err = a.StartNewJob(serviceById, parsedSchedule)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("error starting job"))
		return
	}
	c.Status(http.StatusOK)
}
