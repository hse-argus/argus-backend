package app

import "net/http"

func (a *App) GetAllServices(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}