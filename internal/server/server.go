package server

import (
	"observer/internal/app"
	"observer/internal/config"
	"observer/internal/logger"
	"observer/internal/middleware"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
)

func NewServer(cfg *config.Config, app *app.App) *http.Server {
	router := gin.New()
	router.Use(middleware.EnableCORS)

	authorized := router.Group("/")
	authorized.Use(middleware.JWTTokenVerify)
	{
		authorized.GET("/service", app.GetAllServices)
		authorized.POST("/service", app.AddService)
		authorized.PUT("/service", app.UpdateService)
		authorized.DELETE("/service", app.DeleteService)
		authorized.GET("/service/:id", app.GetServiceById)

		authorized.POST("/healthcheck/:id", app.HealthCheck)
		authorized.POST("/scheduled-healthcheck/:id", app.HandleSchedule)
		authorized.DELETE("/schedule/:id", app.DeleteSchedule)
	}
	router.GET("/ws", app.HandleWSConnection)
	router.POST("/login", app.Login)
	router.POST("/register", app.Register)

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
