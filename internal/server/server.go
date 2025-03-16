package server

import (
	"argus-backend/internal/app"
	"argus-backend/internal/config"
	"argus-backend/internal/logger"
	"argus-backend/internal/middleware"
	"context"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/fx"
	"net/http"
)

func NewServer(cfg *config.Config, app *app.App) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/all-services", middleware.EnableCORS(http.HandlerFunc(app.GetAllServices)))
	mux.Handle("/add-service", middleware.EnableCORS(http.HandlerFunc(app.AddService)))
	mux.Handle("/update-service", middleware.EnableCORS(http.HandlerFunc(app.UpdateService)))
	mux.Handle("/delete-service", middleware.EnableCORS(http.HandlerFunc(app.DeleteService)))
	mux.Handle("/service/", middleware.EnableCORS(http.HandlerFunc(app.GetServiceById)))
	mux.Handle("/healthcheck/", middleware.EnableCORS(http.HandlerFunc(app.HealthCheck)))

	mux.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("swagger/swagger/doc.json")))

	return &http.Server{
		Addr:    cfg.WebAddr,
		Handler: mux,
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
