package server

import (
	"argus-backend/internal/app"
	"argus-backend/internal/config"
	"fmt"
	"net/http"
)

func NewServer(cfg *config.Config, app *app.App) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/all-services", http.HandlerFunc(app.GetAllServices))
	mux.Handle("/add-service", http.HandlerFunc(app.AddService))
	mux.Handle("/update-service", http.HandlerFunc(app.UpdateService))
	mux.Handle("/delete-service", http.HandlerFunc(app.DeleteService))

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
