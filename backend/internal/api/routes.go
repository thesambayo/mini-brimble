package api

import "net/http"

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /deployments", app.handleGetDeployments)
	mux.HandleFunc("POST /deployments", app.handleCreateDeployment)
	mux.HandleFunc("GET /deployments/{id}/logs", app.handleDeploymentLogs)
	return corsMiddleware(mux)
}
