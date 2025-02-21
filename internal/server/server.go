package server

import (
	"argus-backend/internal/app"
	"argus-backend/internal/config"
	"net/http"
)

func NewServer(cfg *config.Config, app *app.App) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/all-services", http.HandlerFunc(app.GetAllServices))
	mux.Handle("/add-servuce", http.HandlerFunc(app.AddService))

	return &http.Server{
		Addr:    cfg.WebAddr,
		Handler: mux,
	}
}
