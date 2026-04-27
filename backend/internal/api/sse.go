package api

import (
	"fmt"
	"net/http"
)

func (app *Application) handleDeploymentLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	deploymentId := r.PathValue("id")

	existingLogs, err := app.Logs.GetByDeploymentID(deploymentId)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	for _, existingLog := range existingLogs {
		fmt.Fprintf(w, "data: %s\n\n", existingLog.Line)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}

	app.StreamsMu.Lock()
	logCh, ok := app.LogStreams[deploymentId]
	app.StreamsMu.Unlock()
	if !ok {
		return
	}

	for {
		select {
		case line, ok := <-logCh:
			if !ok {
				return // channel closed, build finished
			}
			fmt.Fprintf(w, "data: %s\n\n", line)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-r.Context().Done():
			return // client disconnected
		}
	}

}
