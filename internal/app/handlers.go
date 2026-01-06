package app

import (
	"observer/internal/logger"
	"observer/internal/middleware"
	"observer/internal/repository/schedule"
	"observer/internal/repository/service"
	"observer/internal/repository/user"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (a *App) Register(c *gin.Context) {
	logger.Info("/register")

	newUser := user.User{}
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in input data"})
		return
	}

	token, err := a.userService.Register(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (a *App) Login(c *gin.Context) {
	logger.Info("/login")

	newUser := user.User{}
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in input data"})
		return
	}

	token, err := a.userService.Login(newUser.Login, newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (a *App) GetAllServices(c *gin.Context) {
	logger.Info("/get_all_services")

	userId, ok := c.Get("id")
	if !ok {
		c.JSON(400, gin.H{"error": "error getting user id"})
		return
	}

	services, err := a.infoService.GetAllServices(userId.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting all services: " + err.Error()})
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

	userId, ok := c.Get("id")
	if !ok {
		c.JSON(400, gin.H{"error": "error getting user id"})
		return
	}
	newService.UserID = userId.(int)

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

	userId, ok := c.Get("id")
	if !ok {
		c.JSON(400, gin.H{"error": "error getting user id"})
		return
	}

	userIdCasted := strconv.Itoa(userId.(int))
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
		a.connections,
		userIdCasted)
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

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

func (a *App) HandleWSConnection(c *gin.Context) {
	logger.Info("/ws")
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("error connecting to websocket: " + err.Error())
		return
	}

	token := c.Query("token")
	logger.Info("token: " + token)
	claims, err := middleware.ParseToken(token)
	if err != nil {
		logger.Error("error parsing token: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "error parsing token: " + err.Error()})
		return
	}

	logger.Info("id: " + strconv.Itoa(claims.Id))
	a.mu.Lock()
	a.connections[strconv.Itoa(claims.Id)] = ws
	a.mu.Unlock()

	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()
		defer ws.Close()

		for {
			select {
			case <-ticker.C:
				a.mu.Lock()
				err = ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait))
				a.mu.Unlock()

				if err != nil {
					logger.Info("ping failed, closing connection")
					return
				}
			}
		}
	}()

	defer func() {
		a.mu.Lock()
		delete(a.connections, strconv.Itoa(claims.Id))
		a.mu.Unlock()
		logger.Info("ws connection closed")
	}()

	logger.Info("ws connection established")
	for {
		_, _, err = ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Info(fmt.Sprintf("websocket connection closed unexpectedly: %v", err))
			} else {
				logger.Info(fmt.Sprintf("WebSocket connection closed: %v", err))
			}
			break
		}
	}
}

func (a *App) HandleSchedule(c *gin.Context) {
	logger.Info("/health_check-scheduled")

	timeSchedule := schedule.Schedule{}
	if err := c.BindJSON(&timeSchedule); err != nil {
		c.JSON(400, gin.H{"error": "error parsing body"})
		return
	}

	userId, ok := c.Get("id")
	if !ok {
		c.JSON(400, gin.H{"error": "error getting user id"})
		return
	}
	userIdCasted := strconv.Itoa(userId.(int))

	id, _ := strconv.Atoi(c.Param("id"))
	serviceById, err := a.infoService.GetServiceById(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "error getting service"})
		return
	}

	parsedSchedule, ok := timeSchedule.ParseScheduleDuration()
	if !ok {
		c.JSON(400, gin.H{"error": "error parsing schedule duration"})
		return
	}

	jobID, err := a.StartNewJob(serviceById, parsedSchedule, userIdCasted)
	if err != nil {
		c.JSON(500, gin.H{"error": "error starting jo"})
		return
	}

	err = a.infoService.AddJob(id, jobID)
	if err != nil {
		c.JSON(500, gin.H{"error": "error adding job"})
		return
	}

	c.Status(http.StatusOK)
}

func (a *App) DeleteSchedule(c *gin.Context) {
	logger.Info("/delete-schedule")

	id, _ := strconv.Atoi(c.Param("id"))
	jobID, err := a.infoService.DeleteJob(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "error deleting job"})
		return
	}

	err = a.RemoveJob(jobID)
	if err != nil {
		c.JSON(500, gin.H{"error": "error removing job"})
		return
	}

	c.Status(http.StatusOK)
}
