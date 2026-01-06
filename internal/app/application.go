package app

import (
	"observer/internal/logger"
	"observer/internal/repository/service"
	clientservice "observer/internal/service/client-service"
	notificationservice "observer/internal/service/notification-service"
	servicesinfo "observer/internal/service/services-info"
	userservice "observer/internal/service/user-service"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type App struct {
	infoService         servicesinfo.ServicesInfoInterface
	userService         userservice.UserServiceInterface
	clientService       *clientservice.ClientService
	notificationService *notificationservice.NotificationService

	connections map[string]*websocket.Conn
	mu          *sync.RWMutex

	scheduler *gocron.Scheduler
}

func NewApp(infoService servicesinfo.ServicesInfoInterface,
	service *clientservice.ClientService,
	userService userservice.UserServiceInterface) *App {
	scheduler, _ := gocron.NewScheduler()
	return &App{
		infoService:   infoService,
		clientService: service,
		userService:   userService,
		connections:   make(map[string]*websocket.Conn),
		mu:            &sync.RWMutex{},
		scheduler:     &scheduler,
	}
}

func (a *App) HealthCheckTask(service *service.Service, userLogin string) {
	logger.Info(fmt.Sprintf("Health check task for service %s started", service.Name))

	statusCode, err := a.clientService.SendRequest(service.Address, service.Port)
	if err != nil {
		logger.Error("Error sending health check requests" + err.Error())
		return
	}

	logger.Info("Preparing to send web notification")
	err = a.notificationService.SendWebNotification(fmt.Sprintf("Результат healthcheck: %d для сервиса %s",
		statusCode,
		service.Name),
		a.connections,
		userLogin)
	if err != nil {
		logger.Error("Error sending web notification: " + err.Error())
	}

	logger.Info("Preparing to send email notification")
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

func (a *App) StartNewJob(service *service.Service, duration time.Duration, userId string) (uuid.UUID, error) {
	job, err := (*a.scheduler).NewJob(
		gocron.DurationJob(duration),
		gocron.NewTask(
			a.HealthCheckTask, service, userId,
		))
	if err != nil {
		logger.Error("error starting cron job: " + err.Error())
		return uuid.Nil, err
	}
	logger.Info("cron job started")
	return job.ID(), nil
}

func (a *App) RemoveJob(id uuid.UUID) error {
	if err := (*a.scheduler).RemoveJob(id); err != nil {
		logger.Error("error removing cron job: " + err.Error())
		return err
	}

	return nil
}
