package server

import (
	"argus-backend/internal/app"
	"argus-backend/internal/config"
	"argus-backend/internal/middleware"
	"fmt"
	"net/http"
)

func NewServer(cfg *config.Config, app *app.App) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/all-services", middleware.EnableCORS(http.HandlerFunc(app.GetAllServices)))
	mux.Handle("/add-service", middleware.EnableCORS(http.HandlerFunc(app.AddService)))
	mux.Handle("/update-service", middleware.EnableCORS(http.HandlerFunc(app.UpdateService)))
	mux.Handle("/delete-service", middleware.EnableCORS(http.HandlerFunc(app.DeleteService)))
	mux.Handle("/get_service_by_id", middleware.EnableCORS(http.HandlerFunc(app.GetServiceById)))

	return &http.Server{
		Addr:    cfg.WebAddr,
		Handler: mux,
	}
}

func RunServer(srv *http.Server) error {
	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("server failed to start or finished with error: %w", err)
	}

	return nil
}
