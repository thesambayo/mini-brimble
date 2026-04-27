package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"mini-brimble/backend/internal/models"
)

func (app *Application) handleGetDeployments(w http.ResponseWriter, r *http.Request) {
	deployments, err := app.Deployments.GetAll()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"deployments": deployments})
}

func (app *Application) handleCreateDeployment(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SourceType string `json:"source_type"`
		Source     string `json:"source"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if body.SourceType == "" || body.Source == "" {
		http.Error(w, "source_type and source are required", http.StatusBadRequest)
		return
	}

	d := &models.Deployment{
		SourceType: body.SourceType,
		Source:     body.Source,
	}

	if _, err := app.Deployments.Create(d); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	go func() {
		// create associated LogStreams chan
		logChan := make(chan string, 100)
		app.StreamsMu.Lock()
		app.LogStreams[d.Id] = logChan
		app.StreamsMu.Unlock()
		// when done building, or building failed, close the associated LogStreams chan
		defer func() {
			app.StreamsMu.Lock()
			delete(app.LogStreams, d.Id)
			app.StreamsMu.Unlock()
			close(logChan)
		}()

		updateStatus := func(status models.Status, message string) {
			if err := app.Deployments.UpdateStatus(d.Id, status); err != nil {
				log.Printf("failed to update status: %v", err)
			}
			logLine := fmt.Sprintf("[BRIMBLE] %s", message)
			// insert into SQLite
			app.Logs.Insert(&models.Log{
				DeploymentId: d.Id,
				Line:         logLine,
			})
			// send to SSE channel if active
			select {
			case logChan <- logLine:
			default:
			}
		}

		updateStatus(models.StatusBuilding, "Status: building — cloning repository")

		workspacePath, err := app.Pipeline.Clone(d.Id, d.Source)
		if err != nil {
			log.Printf("clone failed for deployment %s: %v", d.Id, err)
			// clone failed
			updateStatus(models.StatusFailed, fmt.Sprintf("Status: failed — clone error: %v", err))
			return
		}

		appName := GetAppNameFromRepo(d.Source)
		imageTag, err := app.Pipeline.Build(d.Id, workspacePath, appName, logChan)
		if err != nil {
			log.Printf("build failed for deployment %s: %v", d.Id, err)
			// build failed
			updateStatus(models.StatusFailed, fmt.Sprintf("Status: failed — build error: %v", err))
			return
		}

		updateStatus(models.StatusDeploying, "Status: deploying — starting container")
		log.Printf("build successful for deployment %s, image: %s", d.Id, imageTag)

		_, err = app.Pipeline.Run(d.Id, imageTag)

		if err != nil {
			updateStatus(models.StatusFailed, fmt.Sprintf("Status: failed — deployment error: %v", err))
			return
		}

		publicUrl, err := app.Pipeline.RegisterRoute(d.Id)
		if err != nil {
			log.Printf("failed to register route: %v", err)
		}

		// deploymentURL := fmt.Sprintf("http://localhost:%d", port)

		if err := app.Deployments.UpdateDeployURL(d.Id, publicUrl); err != nil {
			log.Printf("failed to update deployment url: %v", err)
		}

		if err := app.Deployments.UpdateImageTag(d.Id, imageTag); err != nil {
			log.Printf("failed to update deployment url: %v", err)
		}

		updateStatus(models.StatusRunning, "Status: running — app running")
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(d)
}
