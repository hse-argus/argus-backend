package app

import (
	"argus-backend/internal/logger"
	"argus-backend/internal/repository/service"
	clientservice "argus-backend/internal/service/client-service"
	notificationservice "argus-backend/internal/service/notification-service"
	servicesinfo "argus-backend/internal/service/services-info"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type App struct {
	infoService         servicesinfo.ServicesInfoInterface
	clientService       *clientservice.ClientService
	notificationService *notificationservice.NotificationService

	connections map[*websocket.Conn]bool
	mu          *sync.RWMutex

	scheduler *gocron.Scheduler
}

func NewApp(infoService servicesinfo.ServicesInfoInterface, service *clientservice.ClientService) *App {
	scheduler, _ := gocron.NewScheduler()
	return &App{
		infoService:   infoService,
		clientService: service,
		connections:   make(map[*websocket.Conn]bool),
		mu:            &sync.RWMutex{},
		scheduler:     &scheduler,
	}
}

func (a *App) HealthCheckTask(service *service.Service) {
	logger.Info(fmt.Sprintf("Health check task for service %s started", service.Name))

	statusCode, err := a.clientService.SendRequest(service.Address, service.Port)
	if err != nil {
		logger.Error("Error sending health check requests" + err.Error())
		return
	}

	err = a.notificationService.SendWebNotification(fmt.Sprintf("Результат healthcheck: %d для сервиса %s",
		statusCode,
		service.Name),
		a.connections)
	if err != nil {
		logger.Error("Error sending web notification: " + err.Error())
	}

	err = a.notificationService.SendEmailNotification(fmt.Sprintf("Результат healthcheck: %d для сервиса %s",
		statusCode,
		service.Name))
	if err != nil {
		logger.Error("Error sending email notification: " + err.Error())
	}
}

func InvokeScheduler(s *App) {
	go func() {
		(*s.scheduler).Start()
		logger.Info("Scheduler is started")
		select {}
	}()
}

func (a *App) StartNewJob(service *service.Service, duration time.Duration) error {
	_, err := (*a.scheduler).NewJob(
		gocron.DurationJob(duration),
		gocron.NewTask(
			a.HealthCheckTask, service,
		))
	if err != nil {
		logger.Error("error starting cron job: " + err.Error())
		return err
	}
	logger.Info("cron job started")
	return nil
}
