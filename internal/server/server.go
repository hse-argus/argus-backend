package server

import (
	"argus-backend/internal/app"
	"argus-backend/internal/config"
	"argus-backend/internal/logger"
	"argus-backend/internal/middleware"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
)

func NewServer(cfg *config.Config, app *app.App) *http.Server {
	router := gin.New()
	router.Use(middleware.EnableCORS)

	router.GET("/service", app.GetAllServices)
	router.POST("/service", app.AddService)
	router.PUT("/service", app.UpdateService)
	router.DELETE("/service", app.DeleteService)
	router.GET("/service/:id", app.GetServiceById)

	router.POST("/healthcheck/:id", app.HealthCheck)
	router.POST("/scheduled-healthcheck/:id", app.HandleSchedule)

	router.GET("/ws", app.HandleWSConnection)

	return &http.Server{
		Addr:    cfg.WebAddr,
		Handler: router,
	}
}

func RunServer(lc fx.Lifecycle, srv *http.Server) error {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				logger.Info("starting server on " + srv.Addr)

				if err := srv.ListenAndServe(); err != nil {
					logger.Error("error starting server: " + err.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return nil
}
